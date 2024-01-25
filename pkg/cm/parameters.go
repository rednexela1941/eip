package cm

import (
	"time"

	"github.com/rednexela1941/eip/pkg/cip"
)

//go:generate stringer -type=ConnectionType
type ConnectionType cip.USINT

const (
	// See Volume 1: 3-5.6.1.1.2
	ConnectionTypeNull         ConnectionType = 0b00
	ConnectionTypeMulticast    ConnectionType = 0b01
	ConnectionTypePointToPoint ConnectionType = 0b10
	ConnectionTypeReserved     ConnectionType = 0b11
)

//go:generate stringer -type=ConnectionPriority
type ConnectionPriority cip.USINT

const (
	// See Volume 1: 3-5.6.1.1.1
	ConnectionPriorityLow       ConnectionPriority = 0b00
	ConnectionPriorityHigh      ConnectionPriority = 0b01
	ConnectionPriorityScheduled ConnectionPriority = 0b10
	ConnectionPriorityUrgent    ConnectionPriority = 0b11
)

type (
	// See Volume 1: Table 3-5.11
	// And Volume 1: Table 3-5.12
	// Shared, unpacked structure representing the connection parameters
	Parameters struct {
		RPI            cip.UDINT // microseconds.
		Size           cip.UINT
		Type           ConnectionType
		Priority       ConnectionPriority
		Variable       bool
		RedundantOwner bool
		IsLarge        bool // true if large forward open.
	}
)

// local iface for shared parsing function.
type intoAble interface {
	Into() Parameters
}

func (self *Parameters) GetRPI() time.Duration {
	return time.Duration(self.RPI) * time.Microsecond
}

func (self ConnectionType) IsNull() bool {
	return self == ConnectionTypeNull
}
func (self ConnectionType) IsReserved() bool {
	return self == ConnectionTypeReserved
}
func (self ConnectionType) IsPointToPoint() bool {
	return self == ConnectionTypePointToPoint
}
func (self ConnectionType) IsMulticast() bool {
	return self == ConnectionTypeMulticast
}
func (self *Parameters) IsNull() bool {
	return self.Type.IsNull()
}
func (self *Parameters) IsReserved() bool {
	return self.Type.IsReserved()
}

func (self *_ConnectionParameters) Into() Parameters {
	var p Parameters
	p.IsLarge = false
	p.RPI = self.RPI
	p.Size = cip.UINT(self.Parameters & 0xff)
	p.Variable = 0 != 0b1&(self.Parameters>>9)
	p.Type = ConnectionType(0b11 & (self.Parameters >> 13))
	p.Priority = ConnectionPriority(0b11 & (self.Parameters >> 10))
	p.RedundantOwner = 0 != 0b1&(self.Parameters>>15)
	return p
}

func (self *_LargeConnectionParameters) Into() Parameters {
	var p Parameters
	p.IsLarge = true
	p.RPI = self.RPI
	p.Size = cip.UINT(self.Parameters & 0xffff)
	p.Variable = 0 != 0b1&(self.Parameters>>25)
	p.Type = ConnectionType(0b11 & (self.Parameters >> 29))
	p.Priority = ConnectionPriority(0b11 & (self.Parameters >> 26))
	p.RedundantOwner = 0 != 0b1&(self.Parameters>>31)
	return p
}
