package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

type _SockaddrInfoItemWriter struct {
	ItemHeader
	SockaddrInfo
}

func (self *_Writer) AddSockaddrInfoItem(t ItemID, info SockaddrInfo) {
	self.Items = append(self.Items, &_SockaddrInfoItemWriter{
		ItemHeader: ItemHeader{
			TypeID: t,
		},
		SockaddrInfo: info,
	})
}

func (self *_SockaddrInfoItemWriter) WriteTo(w bbuf.Writer) error {
	temp := bbuf.New(nil)
	temp.Wb(&self.SockaddrInfo)

	if temp.Error() != nil {
		return temp.Error()
	}

	data := temp.Bytes()
	self.Length = cip.UINT(len(data))
	w.Wl(&self.ItemHeader)
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}
