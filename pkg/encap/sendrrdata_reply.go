package encap

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cpf"
)

type SendRRDataReply interface {
	Reply
	cpf.Writer
}

type _SendRRDataReply struct {
	Packet
	DataPacketHeader // these fields are ignored by receiver, so we can leave them zero.
	cpf.Writer
}

func NewSendRRDataReply() SendRRDataReply {
	r := new(_SendRRDataReply)
	r.Writer = cpf.NewWriter()
	return r
}

func (self *_SendRRDataReply) Encode() ([]byte, error) {
	data, err := self.Writer.Encode()
	if err != nil {
		return nil, err
	}
	self.Length = cip.UINT(len(data)) + DataPacketHeaderSize
	buffer := bbuf.New(nil)
	buffer.Wl(self.Header)
	buffer.Wl(self.DataPacketHeader)
	if _, err := buffer.Write(data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), buffer.Error()
}
