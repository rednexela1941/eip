package adapter

import (
	"math/rand"

	"github.com/rednexela1941/eip/pkg/encap"
)

func (self *_Adapter) CanRegisterSession(c *RequestContext) bool {
	// TODO: track sessions by tcp connection.
	if !c.IsTCP() {
		// no sessions on UDP
		return false
	}

	return true
}

func (self *_Adapter) IsValidSession(c *RequestContext, handle encap.SessionHandle) bool {
	return true
}

func (self *_Adapter) RegisterSession(c *RequestContext) encap.SessionHandle {
	// TODO: track sessions by TCP connection.
	self.sessionHandleOffset++
	handle := rand.Uint32() & 0xffff0000 // random 16 bits + incrementing, to make it "readable"
	return handle + self.sessionHandleOffset
}

func (self *_Adapter) _HandleRegisterSession(c *RequestContext, p encap.RegisterSessionRequest) (encap.RegisterSessionReply, error) {
	if p.GetOptions() != 0 {
		// discard non-zero options requests.
		return nil, nil
	}

	reply := encap.NewRegisterSessionReply()
	reply.SetHeader(p.GetHeader())

	if p.GetProtocolVersion() != encap.ProtocolVersion {
		reply.SetStatus(encap.StatusUnsupportedProtocolVersion)
		return reply, nil
	}

	if !self.CanRegisterSession(c) {
		reply.SetStatus(encap.StatusInvalidCommand)
		return reply, nil
	}

	// Also, there is an error code for too many sessions that we could use.

	handle := self.RegisterSession(c)

	reply.SetSessionHandle(handle)

	return reply, nil
}
