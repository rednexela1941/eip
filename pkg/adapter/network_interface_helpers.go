package adapter

import (
	"net"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/network"
)

type (
	// NetworkInterface represents an generic network interface
	// used in the AddNetworkInterface helper.
	NetworkInterface interface {
		GetInterface() *net.Interface

		// Ethernet Link Object releated properties.
		GetSpeed() Mbps
		GetFlags() EthernetLinkInterfaceFlags
		GetType() EthernetLinkType
		GetCapabilities() InterfaceCapabilities

		// TCP/IP Interface Object related properties.
		GetStatus() TCPIPInterfaceStatus
		GetTCPConfigCapability() TCPIPConfigurationCapability
		GetTCPConfigControl() TCPIPConfigurationControl
		GetHostname() string

		GetAddresses() []*network.InterfaceAddr
	}

	// Mbps: cip.UDINT in megabits per second.
	Mbps = cip.UDINT

	// InterfaceCapabilities See Volume 2: 5-5.3.2 Attribute 11
	InterfaceCapabilities struct {
		CapabilityBits   EthernetLinkCapability
		SpeedDuplexArray []SpeedDuplex
	}

	// SpeedDuplex See Volume 2: 5-5.3.2 Attribute 11
	SpeedDuplex struct {
		Speed      cip.UINT
		DuplexMode InterfaceDuplexMode
	}

	MacAddr [6]cip.USINT
)

// GetMACaddr (fixed size)
func GetMACAddr(hw net.HardwareAddr) MacAddr {
	addr := MacAddr{}
	for i, v := range hw {
		if i > 6 {
			break
		}
		addr[i] = v
	}
	return addr
}

func NewInterfaceConfiguration(addr *network.InterfaceAddr) *InterfaceConfiguration {
	c := new(InterfaceConfiguration)
	c.InterfaceAddr = *addr
	return c
}

type portContext struct {
	PortInstanceID       cip.UINT
	EncapsulationTimeout *cip.UINT
}

func newPortContext() *portContext {
	to := new(cip.UINT)
	*to = TCPIPDefaultEncapsulationInvactivityTimeout
	return &portContext{
		EncapsulationTimeout: to,
	}
}

// these are the arguments to listen on.
type listenParams struct {
	addr  *network.InterfaceAddr
	iface *net.Interface
	pctx  *portContext
}

// AddNetworkInterface is a helper that will
// take a custom network interface
// and add appropriate EthernetLink (0xF6) and TCP/IP interface (0xF5)
// objects to the adapter.
func (self *Adapter) AddNetworkInterface(iface NetworkInterface) error {
	ethlinkInstance := self._AddEthernetLinkInstance(iface)
	for _, addr := range iface.GetAddresses() {
		pctx := newPortContext()
		tcpInstance := self._AddTCPIPInterfaceInstance(
			addr,
			iface,
			ethlinkInstance,
			pctx,
		)
		port := self._AddPortInstance(
			iface,
			addr,
			ethlinkInstance,
			tcpInstance,
		)
		pctx.PortInstanceID = port.InstanceID

		self.listenerParams = append(
			self.listenerParams,
			listenParams{
				addr:  addr,
				iface: iface.GetInterface(),
				pctx:  pctx,
			},
		)
		// if err := self.ListenOn(
		// 	addr,
		// 	iface.GetInterface(),
		// 	pctx,
		// ); err != nil {
		// 	return err
		// }
	}
	return nil
}

// See Volume 1: 3-9 "Port Object Class Definition"
func (self *Adapter) _AddPortInstance(
	iface NetworkInterface,
	addr *network.InterfaceAddr,
	ethlinkInstance *Instance,
	tcpInstance *Instance,
) *Instance {
	c := self.GetPortClass()
	i := c.AddInstance(cip.UINT(len(c.Instances) + 1))
	portNumber := i.InstanceID + 1

	i.AddAttribute(1, "PortType", cip.UINTSize).OnGet(
		GetSingle,
		func(res Response) {
			// See Volume 1: 3-9.2.1.1
			res.Wl(cip.UINT(PortTypeEtherNetIP))
		},
	)
	i.AddAttribute(2, "PortNumber", cip.UINTSize).OnGet(
		GetSingle,
		func(res Response) {
			// See Volume 1: 3-9.2.1.2
			res.Wl(portNumber)
		},
	)
	i.AddAttribute(3, "LogicalLinkObject", 0).OnGet(
		GetSingle,
		func(res Response) {
			// See Volume 1: 3-9.2.1.1
			data, err := tcpInstance.GetPath().Encode(2)
			if err != nil {
				self.Logger.Fatal(err)
			}
			res.Wl(cip.UINT(2)) // UINT!
			res.Wl(data)
		},
	)
	i.AddAttribute(4, "PortName", 0).OnGet(
		GetSingle,
		func(res Response) {
			if err := bbuf.WShortString(
				res,
				iface.GetInterface().Name,
				// fmt.Sprintf("%s %d", iface.GetInterface().Name, addrIndex),
			); err != nil {
				self.Logger.Println(err)
			}
		},
	)
	i.AddAttribute(5, "PortTypeName", 0).OnGet(
		GetSingle,
		func(res Response) {
			if err := bbuf.WShortString(res, PortTypeName); err != nil {
				self.Logger.Println(err)
			}
		},
	)
	i.AddAttribute(7, "PortNumberAndNodeAddress", 0).OnGet(
		GetSingle,
		func(res Response) {
			EncodePortSegment(res, portNumber, addr.IP)
		},
	)
	i.AddAttribute(10, "PortRoutingCapabilities", cip.DWORDSize).OnGet(
		GetSingle,
		func(res Response) {
			res.Wl(
				PortRoutingDefault,
			)
		},
	)
	i.AddAttribute(11, "AssociatedCommunicationObjects", 0).OnGet(
		GetSingle,
		func(res Response) {
			res.Wl(cip.USINT(2)) // 2 items.
			// Encode TCP
			data, err := tcpInstance.GetPath().Encode(2)
			if err != nil {
				self.Logger.Println(err)
			}
			res.Wl(cip.USINT(2))
			res.Wl(data)

			data, err = ethlinkInstance.GetPath().Encode(2)
			if err != nil {
				self.Logger.Println(err)
			}
			res.Wl(cip.USINT(2))
			res.Wl(data)
		},
	)

	return i
}

// See Volume 2: 5-5 "Ethernet Link Object"
func (self *Adapter) _AddEthernetLinkInstance(iface NetworkInterface) *Instance {
	ethlink := self.GetEthernetLinkClass()

	i := ethlink.AddInstance(cip.UINT(len(ethlink.Instances) + 1))

	i.AddAttribute(1, "Interface Speed", cip.UDINTSize).OnGet(
		GetFull,
		func(res Response) {
			res.Wl(iface.GetSpeed())
		},
	)
	// See Volume 2: 5-5.3.2.1
	i.AddAttribute(2, "Interface Flags", cip.DWORDSize).OnGet(
		GetFull,
		func(res Response) {
			res.Wl(iface.GetFlags())
		},
	)
	i.AddAttribute(3, "Physical Address", 6*cip.USINTSize).OnGet(
		GetFull,
		func(res Response) {
			res.Wl(GetMACAddr(iface.GetInterface().HardwareAddr))
		},
	)
	// See Volume 2: 5-5.3.2.7
	i.AddAttribute(7, "Interface Type", cip.USINTSize).OnGet(
		GetFull,
		func(res Response) {
			res.Wl(iface.GetType())
		},
	)
	// See Volume 2: 5-5.3.2.11
	i.AddAttribute(11, "Interface Capabilities", 0).OnGet(
		GetFull,
		func(res Response) {
			caps := iface.GetCapabilities()
			res.Wl(caps.CapabilityBits)
			res.Wl(cip.USINT(len(caps.SpeedDuplexArray)))
			for _, sd := range caps.SpeedDuplexArray {
				res.Wl(sd.Speed)
				res.Wl(sd.DuplexMode)
			}
		},
	)

	i.OnService(cip.GetAttributesAll, func(req *Request, res Response) {
		res.SetGeneralStatus(cip.StatusServiceNotSupported)
	})
	i.OnService(cip.GetAttributeList, func(req *Request, res Response) {
		res.SetGeneralStatus(cip.StatusServiceNotSupported)
	})

	return i
}

// See Volume 2: 5-4 "TCP/IP Interface Object"
func (self *Adapter) _AddTCPIPInterfaceInstance(
	addr *network.InterfaceAddr,
	iface NetworkInterface,
	ethlinkInstance *Instance,
	pctx *portContext,
) *Instance {

	c := self.GetTCPIPInterfaceClass()
	i := c.AddInstance(cip.UINT(len(c.Instances) + 1))

	i.AddAttribute(1, "Status", cip.DWORDSize).OnGet(
		GetFull,
		func(res Response) { res.Wl(iface.GetStatus()) },
	)
	i.AddAttribute(2, "Configuration Capability", cip.DWORDSize).OnGet(
		GetFull,
		func(res Response) {
			res.Wl(iface.GetTCPConfigCapability())
		},
	)
	i.AddAttribute(3, "Configuration Control", cip.DWORDSize).OnGet(
		GetFull,
		func(res Response) {
			res.Wl(iface.GetTCPConfigControl())
		},
	).OnSet(
		SetSingle,
		func(req *Request, res Response) {
			res.SetGeneralStatus(cip.StatusDeviceStateConflict)
		},
	)
	i.AddAttribute(4, "Physical Link Object", 0).OnGet(
		GetFull,
		func(res Response) { /* TODO */
			data, err := ethlinkInstance.GetPath().Encode(2)
			if err != nil {
				self.Logger.Println(err)
			}
			res.Wl(cip.UINT(2)) // this is in words, but still a UINT
			res.Wl(data)
		},
	)
	i.AddAttribute(5, "Interface Configuration", 0).OnGet(
		GetFull,
		func(res Response) {
			config := NewInterfaceConfiguration(addr)
			// may have to reverse these.
			res.Wl(config.InterfaceAddr.IP.ToUint())
			res.Wl(config.InterfaceAddr.Netmask.ToUint())

			res.Wl(config.Gateway)
			res.Wl(config.NameServer)
			res.Wl(config.NameServer2)
			res.Wl(config.DomainNameLengthAndPad)
		},
	)
	i.AddAttribute(6, "Hostname", 0).OnGet(
		GetSingle,
		func(res Response) { /* TODO */
			hostname := iface.GetHostname()
			l := cip.UINT(len(hostname))
			res.Wl(l)
			res.Wl([]byte(hostname))
			if l%2 != 0 {
				res.Wl(cip.USINT(0)) // pad byte for even number of bytes.
			}
		},
	)

	i.AddAttribute(13, "Encapsulation Inactivity Timeout", cip.UINTSize).OnGet(
		GetSingle,
		func(res Response) {
			res.Wl(pctx.EncapsulationTimeout)
		},
	).OnSet(
		SetSingle,
		func(req *Request, res Response) {
			r := bbuf.New(req.Request.RequestData)
			var newTimeout cip.UINT
			r.Rl(&newTimeout)
			if newTimeout > 3600 {
				res.SetGeneralStatus(cip.StatusInvalidAttributeValue)
				return
			}
			*pctx.EncapsulationTimeout = newTimeout
		},
	)

	i.OnService(cip.GetAttributesAll, func(req *Request, res Response) {
		res.SetGeneralStatus(cip.StatusServiceNotSupported)
	})
	i.OnService(cip.GetAttributeList, func(req *Request, res Response) {
		res.SetGeneralStatus(cip.StatusServiceNotSupported)
	})

	return i
}
