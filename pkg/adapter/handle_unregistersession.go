package adapter

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/encap"
)

func (self *_Adapter) UnregisterSession(c *RequestContext, handle encap.SessionHandle) {
	// TODO: do cleanup session.
	return
}

func (self *_Adapter) _HandleUnregisterSession(c *RequestContext, p encap.UnregisterSessionRequest) (encap.Reply, error) {
	// TODO: close TCP connection.
	self.UnregisterSession(c, p.GetSessionHandle())
	return nil, fmt.Errorf("unregister session: 0x%X", p.GetSessionHandle())
}
