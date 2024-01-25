// Package identity implements the shared data structure for identity objects (also used in ListIdentity encapsulation command)
package identity

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

type Revision struct {
	Major cip.USINT
	Minor cip.USINT
}

func (self *Revision) String() string {
	return fmt.Sprintf("%d.%d", self.Major, self.Minor)
}

// See Volume 2: Table 2-4.4 (for encapsulation version)
type Identity struct {
	VendorID     cip.UINT
	DeviceType   cip.UINT
	ProductCode  cip.UINT
	Revision     Revision
	Status       cip.WORD
	SerialNumber cip.UDINT
	ProductName  string    // cip.SHORT_STRING
	State        cip.USINT // In list identity, if not implemented this should be 0xFF
}

func (self *Identity) WriteTo(w bbuf.Writer) error {
	w.Wl(&self.VendorID)
	w.Wl(&self.DeviceType)
	w.Wl(&self.ProductCode)
	w.Wl(&self.Revision)
	w.Wl(&self.Status)
	w.Wl(&self.SerialNumber)
	if err := bbuf.WShortString(w, self.ProductName); err != nil {
		return err
	}
	w.Wl(&self.State)
	return w.Error()
}

// func (self *Identity) SetProductName(name string) {
// 	fullLength := len(name)
// 	if fullLength > 0xff {
// 		name = name[:0xff] // chomp.
// 		fullLength = len(name)
// 	}
// 	b := bbuf.New(nil)
// 	b.Wl(cip.USINT(fullLength))
// 	b.Write([]byte(name))
// 	self.ProductName = b.Bytes()
// }
