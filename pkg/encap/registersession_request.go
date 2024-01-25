package encap

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

/*
An originator shall send a RegisterSession command to a target to initiate a session. The
RegisterSession command does not require that a session be established.
NOTE: See section 2-2.1.3, for detailed information on establishing and maintaining a session
*/

// Volume 2: Table 2-4.9
type RegisterSessionRequest interface {
	Request
	GetProtocolVersion() cip.UINT
	GetOptionsFlags() cip.UINT // no options currently defined, shall be zero. discard non-zero
}

type _RegisterSessionRequest struct {
	*Packet
	ProtocolVersion cip.UINT
	OptionsFlags    cip.UINT
}

func (self *Packet) ToRegisterSessionRequest() (RegisterSessionRequest, error) {
	if self.Length != 4 {
		return nil, fmt.Errorf("invalid length (%d) for %s", self.Length, self.Command.String())
	}
	b := bbuf.New(self.CommandSpecificData) // check errors later (length)
	r := new(_RegisterSessionRequest)
	r.Packet = self
	b.Rl(&r.ProtocolVersion)
	b.Rl(&r.OptionsFlags)
	return r, b.Error()
}

func (self *_RegisterSessionRequest) GetProtocolVersion() cip.UINT { return self.ProtocolVersion }
func (self *_RegisterSessionRequest) GetOptionsFlags() cip.UINT    { return self.OptionsFlags }
