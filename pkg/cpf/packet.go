// Package cpf implements the basic Common Packet Format defintions and functionality (Volume 2: 2-5)
package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

// See Volume 2: 2-5.1

type (
	// Common Packet Format (Vol 2: Table 2-5.1)
	Packet struct {
		ItemCount cip.UINT
		Items     []Item
	}

	ItemHeader struct {
		TypeID ItemID
		Length cip.UINT
	}

	// Common Packet Format Item (Vol 2: Table 2-5.2)
	Item struct {
		ItemHeader
		Data []cip.OCTET
	}
)

// read Header.Length bytes from reader.
func (self *ItemHeader) _getData(r bbuf.Reader) ([]byte, error) {
	if self.Length == 0 {
		return nil, nil
	}
	d := make([]cip.OCTET, self.Length)
	r.Rl(&d)
	return d, r.Error()
}
