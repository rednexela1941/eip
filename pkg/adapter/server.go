package adapter

import (
	"fmt"
	"math"
	"net"
	"time"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cpf"
	"github.com/rednexela1941/eip/pkg/encap"
	"github.com/rednexela1941/eip/pkg/network"
	"golang.org/x/net/ipv4"
)

type listenContext struct {
	IOConn *ipv4.PacketConn
	Addr   *network.InterfaceAddr
	portContext
}

// ListenOnIP starts appropriate TCP and UDP servers
// for the given IP address.
func (self *Adapter) ListenOn(
	addr *network.InterfaceAddr,
	iface *net.Interface,
	pctx *portContext,
) error {
	ctx := &listenContext{
		Addr:        addr,
		portContext: *pctx,
	}

	tcpAddr, err := net.ResolveTCPAddr(
		"tcp",
		fmt.Sprintf("%s:%d", addr.IP.String(), network.TCPPort),
	)
	if err != nil {
		return err
	}
	udpAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", addr.IP.String(), network.UDPPort),
	)
	if err != nil {
		return err
	}
	udpIOAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", addr.IP.String(), network.UDPIOPort),
	)
	if err != nil {
		return err
	}

	ioConn, err := self.StartUDPIOServer(udpIOAddr, iface, ctx)
	if err != nil {
		return err
	}
	ctx.IOConn = ioConn

	if err := self.StartTCPServer(tcpAddr, ctx); err != nil {
		return err
	}
	if err := self.StartUDPServer(udpAddr, ctx); err != nil {
		return err
	}
	return nil
}

func (self *Adapter) StartTCPServer(
	addr *net.TCPAddr,
	ctx *listenContext,
) error {
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	self.Logger.Printf("listening on %s://%s\n", addr.Network(), addr.String())
	go self._StartTCPServer(l, ctx)
	return nil
}

func (self *Adapter) _StartTCPServer(
	l *net.TCPListener,
	ctx *listenContext,
) {
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			self.Logger.Println(err)
			return
		}
		go self._HandleTCPConn(c, ctx)
	}
}

func (self *Adapter) _HandleTCPConn(
	c net.Conn,
	ctx *listenContext,
) {
	sessionHandle := encap.SessionHandle(0)
	buffer := make([]byte, math.MaxUint16)
	defer c.Close()

	info, err := network.NewTCPInfo(c, ctx.Addr.Netmask, ctx.IOConn)
	if err != nil {
		self.Logger.Println(err)
		return
	}

	rctx := NewRequestContext(info, ctx.PortInstanceID)

	self.Logger.Println("open", info.String())
	defer func() {
		self.Logger.Println("close", info.String())
	}()

	defer func() {
		if sessionHandle != 0 {
			// remove connection if it exists.
			self.Connections.RemoveConnectionBySessionHandle(sessionHandle)
			self.UnregisterSession(rctx, sessionHandle)
		}
	}()

	for {
		// if it becomes a connection, will want to promote.
		// to something else

		encapTimeout := time.Duration(*ctx.EncapsulationTimeout) * time.Second

		rDeadline := time.Now().Add(encapTimeout)
		if encapTimeout == 0 {
			rDeadline = time.Time{} // zero time.
		}

		if err := c.SetReadDeadline(rDeadline); err != nil {
			self.Logger.Println(err)
			return
		}

		n, err := c.Read(buffer)
		if err != nil {
			self.Logger.Println(err)
			return
		}

		req, err := encap.NewRequest(buffer[:n])
		if err != nil {
			self.Logger.Println(err)
			return
		}
		res, err := self.Handle(rctx, req)
		if err != nil {
			self.Logger.Println(err)
			return
		}
		if res == nil {
			continue
		}
		if req.GetCommand() == encap.RegisterSession {
			// store the session handle so we can clear it out on connection close.
			// by Volume 2: 2-2.1.3.4
			// we need to get rid of sessions when tcp is closed.
			sessionHandle = res.PeekSessionHandle()
		}

		data, err := res.Encode()
		if err != nil {
			self.Logger.Println(err)
			return
		}

		if _, err := c.Write(data); err != nil {
			self.Logger.Println(err)
			return
		}
	}
}

func (self *Adapter) StartUDPIOServer(
	addr *net.UDPAddr,
	iface *net.Interface,
	ctx *listenContext,
) (*ipv4.PacketConn, error) {
	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	p := ipv4.NewPacketConn(l)

	groupAddr := &net.UDPAddr{
		IP: ctx.Addr.GetMulticastAddress().ToIP(),
	}
	if err := p.JoinGroup(iface, groupAddr); err != nil {
		return nil, err
	}
	self.Logger.Printf("listening IO on %s://%s (mcast: %s)\n",
		addr.Network(),
		addr.String(),
		groupAddr.IP.String(),
	)
	go self._StartUDPIOServer(addr, p)
	return p, nil
}

func (self *Adapter) _StartUDPIOServer(localAddr *net.UDPAddr, l *ipv4.PacketConn) {
	defer l.Close()
	buffer := make([]byte, math.MaxUint16)
	for {
		n, _, addr, err := l.ReadFrom(buffer)
		if err != nil {
			self.Logger.Println(err)
			return
		}
		reader, err := cpf.NewIOReader(bbuf.New(buffer[:n]))
		if err != nil {
			self.Logger.Println(err)
			continue
		}
		if reader.GetItemCount() < 2 {
			// not enough items
			continue
		}
		addrItem, ok := reader.GetItem(0).(cpf.SequencedAddressItemReader)
		if !ok {
			self.Logger.Println("not a sequenced address", reader.GetItem(0))
			continue
		}

		dataItem, ok := reader.GetItem(1).(*cpf.Item)
		if !ok {
			self.Logger.Println("invalid data item", dataItem)
			continue
		}
		if dataItem.GetTypeID() != cpf.ConnectedTransportPacket {
			self.Logger.Println("not a connected transport packet", dataItem)
			continue
		}

		conn, ok := self.Connections.GetConnectionByOtoTID(
			addrItem.GetConnectionID(),
		)
		if !ok {
			self.Logger.Println("connection not found", addrItem.GetConnectionID())
			continue
		}

		// TODO: check address is valid for conn.
		noop(addr)

		lastEncapSeq := conn.OtoTEncapsulationSequenceNumber
		currentEncapSeq := addrItem.GetEncapsulationSequenceNumber()
		haveFirstEncapSeqCnt := conn.FirstOtoTEncapsulationSequenceNumberReceived

		if haveFirstEncapSeqCnt && !seqGT32(currentEncapSeq, lastEncapSeq) {
			self.Logger.Printf(
				"encapsulation sequence count invalid last=%d, rxd=%d\n",
				lastEncapSeq,
				currentEncapSeq,
			)
			continue
		}

		conn.FirstOtoTEncapsulationSequenceNumberReceived = true

		cp := conn.Point
		if cp == nil {
			self.Logger.Println("connection point is nil")
			continue
		}

		dataItemReader := bbuf.New(dataItem.Data)
		tc := cp.Transport
		if cp.OutputSize(tc) != dataItem.GetLength() {
			self.Logger.Println(
				"invalid data length",
				cp.OutputSize(tc),
				dataItem.GetLength(),
			)
			continue
		}

		if err := conn.readOtoTIODataHeader(dataItemReader); err != nil {
			self.Logger.Println(err)
			continue
		}
		output := cp.Output
		if output == nil {
			self.Logger.Println("output assembly is nil")
			continue
		}

		if err := output.ReadFrom(dataItemReader); err != nil {
			self.Logger.Println(err)
			continue
		}
		// update sequence counts.
		conn.OtoTEncapsulationSequenceNumber = currentEncapSeq
		conn.UpdateOtoTTimestamp()
	}
}

func (self *Adapter) StartUDPServer(
	addr *net.UDPAddr,
	ctx *listenContext,
) error {
	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	self.Logger.Printf("listening on %s://%s\n", addr.Network(), addr.String())
	go self._StartUDPServer(addr, l, ctx)
	return nil
}

func (self *Adapter) _StartUDPServer(
	localAddr *net.UDPAddr,
	l *net.UDPConn,
	ctx *listenContext,
) {
	defer l.Close()
	buffer := make([]byte, math.MaxUint16)
	for {
		n, addr, err := l.ReadFromUDP(buffer)
		if err != nil {
			self.Logger.Println(err)
			return
		}
		req, err := encap.NewRequest(buffer[:n])
		if err != nil {
			self.Logger.Println(err)
			continue
		}
		info, err := network.NewUDPInfo(
			localAddr,
			addr,
			ctx.Addr.Netmask,
			ctx.IOConn,
		)
		if err != nil {
			self.Logger.Println(err)
			continue
		}
		rctx := NewRequestContext(info, ctx.PortInstanceID)
		res, err := self.Handle(rctx, req)
		if err != nil {
			self.Logger.Println(err)
			continue
		}
		if res == nil {
			// do nothing.
			continue
		}
		data, err := res.Encode()
		if err != nil {
			self.Logger.Println(err)
			continue
		}
		if _, err := l.WriteToUDP(data, addr); err != nil {
			self.Logger.Println(err)
			continue
		}
	}
}
