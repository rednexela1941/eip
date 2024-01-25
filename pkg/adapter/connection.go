package adapter

import (
	"fmt"
	"log"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cm"
	"github.com/rednexela1941/eip/pkg/encap"
	"github.com/rednexela1941/eip/pkg/network"
)

type Connection struct {
	cm.Connection
	network.Info
	Point    *ConnectionPoint
	Channels *ChannelSockaddrs
}

func NewClass0Or1Connection(
	info *network.Info,
	freq *cm.SharedForwardOpenRequest,
	cp *ConnectionPoint,
	channels *ChannelSockaddrs,
) *Connection {
	c := NewConnection(info, freq)
	c.Point = cp
	c.Channels = channels
	c.Info = *info
	return c
}

func NewClass2Or3Connection(
	info *network.Info,
	freq *cm.SharedForwardOpenRequest,
	sessionHandle encap.SessionHandle,
) *Connection {
	c := NewConnection(info, freq)
	c.SessionHandle = sessionHandle
	c.Info = *info
	return c
}

func NewConnection(info *network.Info, freq *cm.SharedForwardOpenRequest) *Connection {
	return &Connection{
		Connection: *cm.NewConnection(freq),
		Info:       *info,
		Point:      nil,
	}
}

// CloseTCP on Class 2/3 connections
func (self *Connection) CloseTCP() error {
	if self.Conn == nil {
		return nil
	}
	return self.Conn.Close()
}

func (self *Connection) readOtoTIODataHeader(r bbuf.Reader) error {
	// will also check sequence count.
	cp := self.Point
	if cp == nil {
		return fmt.Errorf("connection point is nil")
	}

	tc := cp.Transport.TransportClass()
	rtFmt := cp.OtoTFormat

	if tc == cm.Class1 {
		seq := cip.UINT(0)
		r.Rl(&seq)
		if int16(seq-self.OtoTSequenceNumber) > 0 {
			self.OtoTSequenceNumber = seq // update sequence number
		} else {
			return fmt.Errorf("sequence count is less", seq, self.OtoTSequenceNumber)
		}
	}

	if rtFmt == Header32BitFormat {
		h32 := cip.DWORD(0)
		r.Rl(&h32)
		// TODO: something with h32.
	}
	return nil
}

func (self *Connection) writeTtoOIODataHeader(w bbuf.Writer) {
	// write header for class 0/class1 connections
	// See Volume 1: 3-6.1
	cp := self.Point
	if cp == nil {
		log.Fatal("connection point shouldn't be nil")
		return
	}
	tc := cp.Transport.TransportClass()
	fmt := cp.TtoOFormat

	if tc != cm.Class1 && tc != cm.Class0 {
		log.Fatal("invalid transport class %d", tc)
	}

	switch fmt {
	case ModelessFormat, ZeroLengthFormat, HeartbeatFormat:
		if tc == cm.Class1 {
			w.Wl(self.TtoOSequenceNumber)
		}
		return
	case Header32BitFormat:
		// TODO: figure out what to do here.
		h32 := New32BitHeader(true, true, 0)
		if tc == cm.Class0 {
			w.Wl(h32)
		} else if tc == cm.Class1 {
			w.Wl(self.TtoOSequenceNumber)
			w.Wl(h32)
		}
		return
	default:
		log.Fatal("invalid format: %s", fmt.String())
	}
	log.Fatal("unhandled")
}

// See Volume 1: 3-6.1.4 "32-Bit Header Format"
// And Volume 1: 3-6.5.4
// And Volume 1: 3-6.5.4.3.2 "Claim Output Ownership (COO) Flag"
// And Volume 1: 3-6.5.4.3.3 "Ready for Ownership of Outputs (ROO) Priority Value"
func New32BitHeader(running bool, claimOutputOwnership bool, readyForOwnership uint8) cip.DWORD {
	header := cip.DWORD(0)
	// ROO and COO
	if running {
		header |= 0b1 // first bit is run/idle.
	}
	if claimOutputOwnership {
		header |= (0b1 << 1)
		header |= (cip.UDINT(readyForOwnership) & 0b11) << 2
	}
	return header
}

func (self *Connection) String() string {
	s := fmt.Sprintf("Connection(%X, %X){Class:%d",
		self.OtoTNetworkConnectionID,
		self.TtoONetworkConnectionID,
		self.TransportClassAndTrigger.TransportClass(),
	)
	if self.Channels != nil {
		s += ", "
		s += self.Channels.String()
		s += ", " + self.OtoTParameters.GetRPI().String()
		s += ", " + self.TtoOParameters.GetRPI().String()
	}
	s += "}"
	return s
}
