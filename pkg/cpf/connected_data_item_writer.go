package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/mr"
)

type ConnectedDataItemWriter interface {
	IHeaderWriter
	mr.ResponseWriter
}

const ConnectedDataSequenceCountSize cip.UINT = 2

type _ConnectedDataItemWriter struct {
	_UnconnectedDataItemWriter
	SequenceCount cip.UINT
}

func (self *_Writer) AddConnectedDataItem(sequenceCount cip.UINT) ConnectedDataItemWriter {
	item := new(_ConnectedDataItemWriter)
	item.TypeID = ConnectedTransportPacket
	item.SequenceCount = sequenceCount
	item.ResponseWriter = mr.NewResponseWriter()
	self.Items = append(self.Items, item)
	return item
}

func (self *_ConnectedDataItemWriter) WriteTo(w bbuf.Writer) error {
	data, err := self.ResponseWriter.Encode()
	if err != nil {
		return err
	}
	self.Length = cip.UINT(len(data)) + ConnectedDataSequenceCountSize
	w.Wl(self.ItemHeader)
	w.Wl(self.SequenceCount)
	if _, err := w.Write(data); err != nil {
		return err
	}
	return w.Error()
}
