package encap

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cpf"
)

// Volume 2: Table 2-4.3
type ListIdentityReply interface {
	Reply
	cpf.Writer
}

type _ListIdentityReply struct {
	Packet
	cpf.Writer
}

func NewListIdentityReply() ListIdentityReply {
	r := new(_ListIdentityReply)
	r.Writer = cpf.NewWriter()
	return r
}

func (self *_ListIdentityReply) Encode() ([]byte, error) {
	buffer := bbuf.New(nil)
	data, err := self.Writer.Encode()
	if err != nil {
		return nil, err
	}
	self.Length = cip.UINT(len(data)) // update length before sending.
	buffer.Wl(self.Header)
	if _, err := buffer.Write(data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), buffer.Error()
}
