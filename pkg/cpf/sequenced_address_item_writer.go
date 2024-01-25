package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

type _SequencedAddressItemWriter _SequencedAddressItemReader

func (self *_Writer) AddSequencedAddressItem(
	connID cip.UDINT,
	seqNo cip.UDINT,
) {
	self.Items = append(self.Items, &_SequencedAddressItemWriter{
		_ConnectedAddressItemReader: _ConnectedAddressItemReader{
			Item: Item{
				ItemHeader: ItemHeader{
					TypeID: SequencedAddress,
					Length: 8,
				},
			},
			ConnectionID: connID,
		},
		EncapsulatedSequenceNumber: seqNo,
	})
}

func (self *_SequencedAddressItemWriter) WriteTo(w bbuf.Writer) error {
	self.Length = 8
	w.Wl(self.ItemHeader)
	w.Wl(self.ConnectionID)
	w.Wl(self.EncapsulatedSequenceNumber)
	return w.Error()
}
