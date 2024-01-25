package cpf

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

// Volume 2: 2-5.2.2
type ConnectedAddressItemReader interface {
	IReader
	GetConnectionID() cip.UDINT
}

type _ConnectedAddressItemReader struct {
	Item
	ConnectionID cip.UDINT
}

func (self *_ConnectedAddressItemReader) GetConnectionID() cip.UDINT {
	return self.ConnectionID
}

func (self *Item) ToConnectedAddressItemReader(r bbuf.Reader) (ConnectedAddressItemReader, error) {
	if self.Length != 4 {
		return nil, fmt.Errorf("invalid length %d for %s", self.Length, self.TypeID.String())
	}
	item := &_ConnectedAddressItemReader{Item: *self}
	r.Rl(&item.ConnectionID)
	return item, r.Error()
}
