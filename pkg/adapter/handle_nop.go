package adapter

import (
	"github.com/rednexela1941/eip/pkg/encap"
)

// Vol 2: 2-4.1: no reply needed.
func (self *_Adapter) _HandleNOP(_ *RequestContext, _ encap.NOPRequest) (encap.Reply, error) {
	return nil, nil
}
