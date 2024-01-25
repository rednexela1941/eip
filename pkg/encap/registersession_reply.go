package encap

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

// Volume 2: Table 2-4.10
type RegisterSessionReply interface {
	Reply
	SetHighestSupportedProtocolVersion(cip.UINT)
	SetOptionsFlags(cip.UINT)
}

type _RegisterSessionReply struct {
	Packet
	HighestSupportedVersion cip.UINT
	OptionsFlags            cip.UINT // 0
}

func (self *_RegisterSessionReply) SetHighestSupportedProtocolVersion(version cip.UINT) {
	self.HighestSupportedVersion = version
}

func (self *_RegisterSessionReply) SetOptionsFlags(flags cip.UINT) {
	self.OptionsFlags = 0
}

func NewRegisterSessionReply() RegisterSessionReply {
	r := new(_RegisterSessionReply)
	r.HighestSupportedVersion = ProtocolVersion
	return r
}

func (self *_RegisterSessionReply) Encode() ([]byte, error) {
	self.Length = 4
	buffer := bbuf.New(nil)
	buffer.Wl(self.Header)
	buffer.Wl(self.HighestSupportedVersion)
	buffer.Wl(self.OptionsFlags)
	return buffer.Bytes(), buffer.Error()
}
