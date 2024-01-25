package encap

import (
	"github.com/rednexela1941/eip/pkg/cpf"
)

// The specification says that a "reply shall not be returned" (Volume 2: 2-4.8).
// But this is not really true, for class 2/3 PDU (Volume 1)
// See Volume 2: 3-2.3
// And Volume 2: Table 3-2.4
type SendUnitDataReply SendRRDataReply

type _SendUnitDataReply _SendRRDataReply

func NewSendUnitDataReply() SendUnitDataReply {
	r := new(_SendUnitDataReply)
	r.Writer = cpf.NewWriter()
	return r
}

func (self *_SendUnitDataReply) Encode() ([]byte, error) {
	r := (*_SendRRDataReply)(self)
	return r.Encode()
}
