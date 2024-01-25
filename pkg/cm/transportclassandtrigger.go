package cm

import "github.com/rednexela1941/eip/pkg/cip"

// Volume 1: 3-4.5.3
type (
	TransportClassAndTrigger cip.BYTE
)

// Tranport Class
const (
	Class0 TransportClassAndTrigger = 0
	Class1 TransportClassAndTrigger = 1
	Class2 TransportClassAndTrigger = 2
	Class3 TransportClassAndTrigger = 3

	ClassMask TransportClassAndTrigger = 0b1111
)

// Production Trigger
// See Volume 1: Table 3-4.14 "Possible values within Production Trigger Bits"
const (
	// Cyclic - on timer.
	Cyclic TransportClassAndTrigger = 0 << 4
	// ChangeOfState - when data changes.
	ChangeOfState TransportClassAndTrigger = 1 << 4
	// ApplicationObject - whenever we want.
	ApplicationObject TransportClassAndTrigger = 2 << 4

	ProductionTriggerMask TransportClassAndTrigger = 0b111 << 4
)

// Direction
const (
	DirectionClient TransportClassAndTrigger = 0 << 7
	DirectionServer TransportClassAndTrigger = 1 << 7

	DirectionMask TransportClassAndTrigger = 0b1 << 7
)

func (self TransportClassAndTrigger) Direction() TransportClassAndTrigger {
	return self & DirectionMask
}

func (self TransportClassAndTrigger) ProductionTrigger() TransportClassAndTrigger {
	return self & ProductionTriggerMask
}

func (self TransportClassAndTrigger) TransportClass() TransportClassAndTrigger {
	return self & ClassMask
}

func (self TransportClassAndTrigger) TransportClassValid() bool {
	switch self.TransportClass() {
	case Class0, Class1, Class2, Class3:
		return true
	default:
		return false
	}
}

func (self TransportClassAndTrigger) ProductionTriggerValid() bool {
	switch self.ProductionTrigger() {
	case Cyclic, ChangeOfState, ApplicationObject:
		return true
	default:
		return false
	}
}
