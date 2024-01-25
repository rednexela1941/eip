package epath

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

//go:generate stringer -type=Format

// Format specifies the Electronic Key format as defined in
// Volume 1 Appendix C.
type Format cip.USINT

const (
	// See Volume 1: Table C-1.5
	Format4 Format = 4
	// See Volume 1: Table C-1.6
	Format5 Format = 5
)

type MajorRevisionCompatibility cip.BYTE

func (self MajorRevisionCompatibility) GetCompatibilityBit() bool {
	// Table C-1.5
	// when this bit is set, the key values are being used to identify the
	// device receptient is being requested to emulate.
	// 0 is invalid value for fields.
	return (self>>7)&1 != 0
}

func (self MajorRevisionCompatibility) GetRevision() cip.USINT {
	return cip.USINT(self & 0b1111111)
}

type (
	// See Volume 1: Table C-1.5
	ElectronicKeyV4 struct {
		Format        Format
		VendorID      cip.UINT
		DeviceType    cip.UINT
		ProductCode   cip.UINT
		MajorRevision MajorRevisionCompatibility
		MinorRevision cip.USINT
	}

	// See Volume 1: Table C-1.5
	ElectronicKey struct {
		ElectronicKeyV4
		SerialNumber cip.UDINT // only included in table Format5 keys.
	}
)

func (self *ElectronicKey) GetMajorRevision() cip.USINT {
	return self.MajorRevision.GetRevision()
}

func (self *ElectronicKey) GetMinorRevision() cip.USINT {
	return self.MinorRevision
}

func (self *ElectronicKey) ValidFormat() bool {
	switch self.Format {
	case Format4, Format5:
		return true
	default:
		return false
	}
}

// When this is set, ekey is used to identify device and zeroes are invalid values.
func (self *ElectronicKey) GetCompatibilityBit() bool {
	return self.MajorRevision.GetCompatibilityBit()
}

func NewElectronicKey(r bbuf.Reader) (*ElectronicKey, error) {
	k := new(ElectronicKey)
	r.Rl(&k.ElectronicKeyV4)
	if k.Format == Format5 {
		r.Rl(&k.SerialNumber)
	}
	return k, r.Error()
}
