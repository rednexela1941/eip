package encap

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cpf"
)

type (
	// Volume 2: 2-4.6.3
	ListServicesReply ListIdentityReply

	_ListServicesReply _ListIdentityReply
)

func NewListServicesReply() ListServicesReply {
	r := new(_ListServicesReply)
	r.Writer = cpf.NewWriter()
	return r
}

func (self *_ListServicesReply) Encode() ([]byte, error) {
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
