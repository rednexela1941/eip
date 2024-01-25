package encap

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cpf"
)

type (
	// Originator (Scanner) to Target (Adapter) Encapsulation Request.
	Request interface {
		GetHeader() *Header
		GetCommand() Command
		GetLength() cip.UINT
		GetSessionHandle() SessionHandle
		GetSenderContext() SenderContext
		GetOptions() cip.UDINT
	}
)

func (self *Header) GetHeader() *Header              { return self }
func (self *Header) GetCommand() Command             { return self.Command }
func (self *Header) GetLength() cip.UINT             { return self.Length }
func (self *Header) GetSessionHandle() SessionHandle { return self.SessionHandle }
func (self *Header) GetSenderContext() SenderContext { return self.SenderContext }
func (self *Header) GetOptions() cip.UDINT           { return self.Options }

// Create a new Request interface from incoming data. Can be typecast to specific encapsulation command interfaces <Command>Request. The excessive use of interfaces here is primarily for debugging purposes.
func NewRequest(data []byte) (Request, error) {
	p := new(Packet)
	buffer := bbuf.New(data)
	buffer.Rl(&p.Header)
	p.CommandSpecificData = buffer.Bytes()
	if len(p.CommandSpecificData) != int(p.Length) {
		return nil, fmt.Errorf("size mismatch h=%d, len=%d", p.Length, len(p.CommandSpecificData))
	}
	if buffer.Error() != nil {
		return nil, buffer.Error()
	}
	switch p.Command {
	case NOP:
		return p.ToNOPRequest()
	case ListServices:
		return p.ToListServicesRequest()
	case ListIdentity:
		return p.ToListIdentityRequest()
	case ListInterfaces:
		return p.ToListInterfacesRequest()
	case RegisterSession:
		return p.ToRegisterSessionRequest()
	case UnRegisterSession:
		return p.ToUnregisterSessionRequest()
	case SendRRData:
		return p.ToSendRRDataRequest()
	case SendUnitData:
		return p.ToSendUnitDataRequest()
	default:
		return nil, fmt.Errorf("invalid/unsupported encapsulation command: %s", p.Command.String())
	}
}

func GetItemReader(r Request) (cpf.Reader, error) {
	switch req := r.(type) {
	case SendRRDataRequest:
		return req.GetEncapsulatedPacket(), nil
	case SendUnitDataRequest:
		return req.GetEncapsulatedPacket(), nil
	default:
		return nil, fmt.Errorf("cannot get item reader")
	}
}

func GetSockaddrInfoReader(r Request, t cpf.ItemID) (cpf.SockaddrInfoItemReader, bool) {
	if t != cpf.SockaddrInfoOtoT && t != cpf.SockaddrInfoTtoO {
		return nil, false
	}
	reader, err := GetItemReader(r)
	if err != nil {
		return nil, false
	}
	item, ok := reader.GetItemWithType(t)
	if !ok {
		return nil, false
	}
	sockItem, ok := item.(cpf.SockaddrInfoItemReader)
	if !ok {
		return nil, false
	}
	return sockItem, true
}
