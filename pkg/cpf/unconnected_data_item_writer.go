package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/mr"
)

type UnconnectedDataItemWriter interface {
	IHeaderWriter
	mr.ResponseWriter
}

type _UnconnectedDataItemWriter struct {
	ItemHeader
	mr.ResponseWriter
}

func (self *_Writer) AddUnconnectedDataItem() UnconnectedDataItemWriter {
	item := &_UnconnectedDataItemWriter{
		ItemHeader: ItemHeader{
			TypeID: UnconnectedMessage,
		},
		ResponseWriter: mr.NewResponseWriter(),
	}
	self.Items = append(self.Items, item)
	return item
}

func (self *_UnconnectedDataItemWriter) WriteTo(w bbuf.Writer) error {
	data, err := self.ResponseWriter.Encode()
	if err != nil {
		return err
	}
	self.Length = cip.UINT(len(data))
	w.Wl(self.ItemHeader)
	if _, err := w.Write(data); err != nil {
		return err
	}
	return w.Error()
}
