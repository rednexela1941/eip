package cm

import (
	"time"

	"github.com/rednexela1941/eip/pkg/cip"
)

// See Volume 1: 3-5.6.1.2.1
// And Volume 1: Table 3-5.13
type PriorityTimeTick cip.BYTE

func (self PriorityTimeTick) PriorityReserved() bool {
	return 0b1&(self>>4) != 0
}

func (self PriorityTimeTick) PriorityNormal() bool {
	return !self.PriorityReserved()
}

func (self PriorityTimeTick) TimeTickBits() uint8 {
	return uint8(self & 0b1111)
}

// By Volume 1: Table 3-5.14
// this is 2 ^ time_tick milliseconds.
func (self PriorityTimeTick) TickTime() time.Duration {
	exponent := int(self.TimeTickBits())
	return (1 << exponent) * time.Millisecond
}
