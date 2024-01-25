package encap

import "github.com/rednexela1941/eip/pkg/cpf"

// The SendUnitData command shall send encapsulated connected messages.  A reply shall not be
// returned.
// NOTE: When used to encapsulate the CIP, the SendUnitData command is used to send CIP
// connected data in both the OtoT and TtoO directions
// Volume 2: Table 2-4.17
type SendUnitDataRequest SendRRDataRequest

type _SendUnitDataRequest _SendRRDataRequest

func (self *_SendUnitDataRequest) GetEncapsulatedPacket() cpf.Reader { return self.EncapsulatedPacket }

func (self *Packet) ToSendUnitDataRequest() (SendUnitDataRequest, error) {
	r, err := self.ToSendRRDataRequest()
	return (*_SendUnitDataRequest)(r), err
}
