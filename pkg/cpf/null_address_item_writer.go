package cpf

import "github.com/rednexela1941/eip/pkg/bbuf"

type NullAddressItemWriter Item

func (self *_Writer) AddNullAddressItem() {
	self.Items = append(self.Items, &NullAddressItemWriter{
		ItemHeader: ItemHeader{TypeID: AddressNull},
	})
}

func (self *NullAddressItemWriter) WriteTo(w bbuf.Writer) error {
	self.Length = 0
	w.Wl(self.ItemHeader)
	return w.Error()
}
