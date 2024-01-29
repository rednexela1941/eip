package cm

import (
	"fmt"
	"log"
	"time"

	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/encap"
)

// Connection represents an active connection.

type ConnectionTimestamp struct {
	Time  time.Time
	First bool
}

type Connection struct {
	SharedForwardOpenRequest
	SessionHandle                                encap.SessionHandle // optional, for class 3 conns
	OtoTTimestamp                                ConnectionTimestamp
	TtoOTimestamp                                ConnectionTimestamp
	OtoTSequenceNumber                           cip.UINT
	FirstOtoTEncapsulationSequenceNumberReceived bool
	OtoTEncapsulationSequenceNumber              cip.UDINT
	LastSentPacket                               encap.SendUnitDataReply
	TtoOSequenceNumber                           cip.UINT
	TtoOEncapsulationSequenceNumber              cip.UDINT
}

func (self *Connection) String() string {
	return fmt.Sprintf(
		"Connection(serial=%X,ovendor=%X,oserial=%X)",
		self.Triad.ConnectionSerialNumber,
		self.Triad.OriginatorVendorID,
		self.Triad.OriginatorSerialNumber,
	)
}

func NewTimestamp() ConnectionTimestamp {
	return ConnectionTimestamp{
		Time:  time.Now(),
		First: true,
	}
}

func (self *Connection) IsIO() bool {
	tc := self.TransportClassAndTrigger.TransportClass()
	switch tc {
	case Class0, Class1:
		return true
	default:
		return false
	}
}

func (self *ConnectionTimestamp) Update() {
	self.Time = time.Now()
	self.First = false
}

func NewConnection(req *SharedForwardOpenRequest) *Connection {
	// See Volume 1: Table 3-5.16 for Encoded Application Path Ordering
	return &Connection{
		SharedForwardOpenRequest: *req,
		OtoTTimestamp:            NewTimestamp(),
		TtoOTimestamp:            NewTimestamp(),
	}
}

func (self *Connection) UpdateOtoTTimestamp() {
	self.OtoTTimestamp.Update()
}

func (self *Connection) UpdateTtoOTimestamp() {
	self.TtoOTimestamp.Update()
}

func (self *Connection) GetOtoTRPI() time.Duration {
	return time.Duration(self.OtoTParameters.RPI) * time.Microsecond
}

func (self *Connection) GetOtoTTimeout() time.Duration {
	rpi := self.GetTtoORPI()
	return self._GetTimeout(rpi)
}

func (self *Connection) GetTtoORPI() time.Duration {
	return time.Duration(self.TtoOParameters.RPI) * time.Microsecond
}

func (self *Connection) GetTtoOTimeout() time.Duration {
	// multiplier, ok := self.GetTimeoutMultiplier()
	rpi := self.GetTtoORPI()
	return self._GetTimeout(rpi)
}

func (self *Connection) IsTimeToSend(tick time.Duration) bool {
	rpi := self.GetTtoORPI()
	return time.Now().Sub(self.TtoOTimestamp.Time) >= (rpi - tick)
}

func (self *Connection) IsOtoTTimedOut() bool {
	t := self.GetOtoTTimeout()
	if self.OtoTTimestamp.First && t < 10*time.Second {
		// See Volume 1: 3-4.6.2 Inactivity/Watchdow.
		// default timeout of 10 seconds for first TX.
		t = 10 * time.Second
	}
	if time.Now().Sub(self.OtoTTimestamp.Time) > t {
		return true
	}
	return false
}

func (self *Connection) _GetTimeout(rpi time.Duration) time.Duration {
	multiplier, ok := self.GetTimeoutMultiplier()
	if !ok {
		log.Println("invalid timeout multiplier")
		return rpi
	}
	return time.Duration(multiplier) * rpi
}
