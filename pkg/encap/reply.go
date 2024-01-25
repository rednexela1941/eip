package encap

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cpf"
)

type (
	// Target (Adapter) to Originator (Scanner) Encapsulation Reply.
	Reply interface {
		GetCommand() Command // for logging.

		SetHeader(*Header)
		SetStatus(ErrorCode)
		SetSessionHandle(SessionHandle)
		PeekSessionHandle() SessionHandle

		Encode() ([]byte, error)
	}

	// _Reply implements a default encapsulation reply
	_Reply struct {
		Packet
		*bbuf.Buffer
	}
)

func (self *Header) SetStatus(code ErrorCode)              { self.Status = code }
func (self *Header) SetHeader(h *Header)                   { *self = *h }
func (self *Header) SetSessionHandle(handle SessionHandle) { self.SessionHandle = handle }
func (self *Header) PeekSessionHandle() SessionHandle      { return self.SessionHandle }

func _NewReply() *_Reply {
	r := new(_Reply)
	r.Buffer = bbuf.New(nil)
	return r
}

func (self *_Reply) Encode() ([]byte, error) {
	self.Length = cip.UINT(self.Len())
	buffer := bbuf.New(nil)
	buffer.Wl(self.Header)
	buffer.Wl(self.Bytes())
	return buffer.Bytes(), buffer.Error()
}

func GetItemWriter(r Reply) (cpf.Writer, error) {
	switch res := r.(type) {
	case SendRRDataReply:
		return res, nil
	case SendUnitDataReply:
		return res, nil
	default:
		return nil, fmt.Errorf("cannot get item writer")
	}
}
