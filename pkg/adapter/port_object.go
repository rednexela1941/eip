package adapter

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/network"
)

const PortObjectRevision cip.UINT = 2

const PortTypeEtherNetIP cip.UINT = 4
const PortTypeName = "EtherNet/IP"

type PortRoutingCapabilities cip.DWORD

const (
	PortRoutingInUnconnected  PortRoutingCapabilities = 1 << 0
	PortRoutingOutUnconnected PortRoutingCapabilities = 1 << 1
	PortRoutingInClass0And1   PortRoutingCapabilities = 1 << 2
	PortRoutingOutClass0And1  PortRoutingCapabilities = 1 << 3
	PortRoutingInClass2And3   PortRoutingCapabilities = 1 << 4
	PortRoutingOutClass2And3  PortRoutingCapabilities = 1 << 5

	PortRoutingDefault PortRoutingCapabilities = 0b11111
)

func (self *Adapter) GetPortClass() *Class {
	c, ok := self.Classes[cip.PortClassCode]
	if ok {
		return c
	}
	c = self.AddClass("Port", cip.PortClassCode, PortObjectRevision)

	c.AddAttribute(1, "ClassRevision", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) {
			res.Wl(PortObjectRevision)
		},
	)
	c.AddAttribute(2, "MaxInstance", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestInstanceID()) },
	)
	c.AddAttribute(3, "NumInstances", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.NumberOfInstances()) },
	)
	c.AddAttribute(6, "MaxClassAttributeID", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestAttributeID()) },
	)
	c.AddAttribute(7, "MaxInstanceAttributeID", cip.UDINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestInstanceAttributeID()) },
	)
	c.AddAttribute(8, "EntryPort", cip.UINTSize).OnService(
		cip.GetAttributeSingle,
		func(req *Request, res Response) {
			res.Wl(req.PortInstanceID)
		},
	)
	c.AddAttribute(9, "PortInstanceInfo", 0).OnGet(
		GetSingle,
		func(res Response) {
			// instance 0
			res.Wl(cip.UINT(0))
			res.Wl(cip.UINT(0))

			for _, inst := range c.Instances {
				t, ok := inst.Attributes[1] // Port Type
				if !ok {
					res.SetGeneralStatus(cip.StatusAttributeNotSupported)
					return
				}
				n, ok := inst.Attributes[2] // Port Number
				if !ok {
					res.SetGeneralStatus(cip.StatusAttributeNotSupported)
					return
				}
				if err := t.callGetSingle(nil, res); err != nil {
					res.SetGeneralStatus(cip.StatusAttributeNotSupported)
					return
				}
				if err := n.callGetSingle(nil, res); err != nil {
					res.SetGeneralStatus(cip.StatusAttributeNotSupported)
					return
				}
			}
		},
	)

	c.OnService(cip.GetAttributesAll, func(req *Request, res Response) {
		res.SetGeneralStatus(cip.StatusServiceNotSupported)
	})

	return c
}

func EncodePortSegment(w bbuf.Writer, portNum cip.UINT, ip network.IPv4) {
	first := cip.BYTE(0)

	ipStr := ip.ToIP().String()
	first |= 1 << 4 // Extended Link Address Size
	ipLen := cip.USINT(len(ipStr))

	totalLen := 1 + 1 + int(ipLen)
	if portNum >= 0b1111 {
		first |= 0b1111 // fill out to 15
		totalLen += 2
	} else {
		first |= cip.USINT(portNum)
	}
	needPad := totalLen%2 != 0
	if needPad {
		totalLen++
	}
	// w.Wl(cip.USINT(totalLen / 2))
	w.Wl(first)
	w.Wl(ipLen) // Extended Link Size
	if portNum >= 0b1111 {
		w.Wl(portNum) // include large port number.
	}
	w.Wl([]byte(ipStr))
	if needPad {
		w.Wl(cip.BYTE(0))
	}

	return
}
