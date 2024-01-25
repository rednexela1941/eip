package encap

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cpf"
)

/*
A SendRRData command shall transfer an encapsulated request/reply packet between the
originator and target, where the originator initiates the command.  The actual request/reply
packets shall be encapsulated in the data portion of the message and shall be the responsibility
of the target and originator.
NOTE: When used to encapsulate the CIP, the SendRRData request and response are used to
send encapsulated UCMM messages (unconnected messages).  See chapter 3 for more detail.
*/

// Volume 2: Table 2-4.15
type SendRRDataRequest interface {
	Request
	// The Interface handle shall identify the Communications Interface to which the request is directed.  This handle shall be 0 for encapsulating CIP packets.
	GetInterfaceHandle() cip.UDINT // 0
	// The target shall abort the requested operation after the timeout expires.  When the “timeout” field is in the range 1 to 65535, the timeout shall be set to this number of seconds.  When the “timeout” field is set to 0, the encapsulation protocol shall not have its own timeout.  Instead, it shall rely on the timeout mechanism of the encapsulated protocol.
	// When the SendRRData command is used to encapsulate CIP packets, the Timeout field shall be set to 0, and shall be ignored by the target.
	GetTimeout() cip.UINT

	GetEncapsulatedPacket() cpf.Reader
}

// Volume 2: Table 2-4.17
const DataPacketHeaderSize cip.UINT = 6

type DataPacketHeader struct {
	InterfaceHandle cip.UDINT
	Timeout         cip.UINT
}

type _SendRRDataRequest struct {
	*Packet
	DataPacketHeader
	EncapsulatedPacket cpf.Reader
	// Todo EncapsulatedPacket
}

func (self *DataPacketHeader) GetInterfaceHandle() cip.UDINT      { return self.InterfaceHandle }
func (self *DataPacketHeader) GetTimeout() cip.UINT               { return self.Timeout }
func (self *_SendRRDataRequest) GetEncapsulatedPacket() cpf.Reader { return self.EncapsulatedPacket }

func (self *Packet) ToSendRRDataRequest() (*_SendRRDataRequest, error) {
	b := bbuf.New(self.CommandSpecificData)
	r := new(_SendRRDataRequest)
	r.Packet = self
	b.Rl(&r.DataPacketHeader)
	epr, err := cpf.NewReader(b)
	if err != nil {
		return r, err
	}
	r.EncapsulatedPacket = epr
	return r, b.Error()
}

// Circular Import Problem?
// cpf.Reader.Parent() -> encap.Request()
// third, shared interface.
