package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

type _ConnectedAddressItemWriter _ConnectedAddressItemReader

func (self *_Writer) AddConnectedAddressItem(connID cip.UDINT) {
	self.Items = append(self.Items, &_ConnectedAddressItemWriter{
		Item: Item{
			ItemHeader: ItemHeader{
				TypeID: AddressConnected,
				Length: 4,
			},
		},
		ConnectionID: connID,
	})
}

func (self *_ConnectedAddressItemWriter) WriteTo(w bbuf.Writer) error {
	self.Length = 4
	w.Wl(self.ItemHeader)
	w.Wl(self.ConnectionID)
	return w.Error()
}
