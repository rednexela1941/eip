package cpf

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

// Volume 2: 2-5.2.3
type SequencedAddressItemReader interface {
	ConnectedAddressItemReader
	GetEncapsulationSequenceNumber() cip.UDINT
}

type _SequencedAddressItemReader struct {
	_ConnectedAddressItemReader
	EncapsulatedSequenceNumber cip.UDINT
}

func (self *_SequencedAddressItemReader) GetConnectionID() cip.UDINT {
	return self.ConnectionID
}

func (self *_SequencedAddressItemReader) GetEncapsulationSequenceNumber() cip.UDINT {
	return self.EncapsulatedSequenceNumber
}

func (self *Item) ToSequencedAddressItemReader(r bbuf.Reader) (SequencedAddressItemReader, error) {
	if self.Length != 8 {
		return nil, fmt.Errorf("invalid length %d for %s", self.Length, self.TypeID.String())
	}
	item := &_SequencedAddressItemReader{}
	item.Item = *self
	r.Rl(&item.ConnectionID)
	r.Rl(&item.EncapsulatedSequenceNumber)
	return item, r.Error()
}
