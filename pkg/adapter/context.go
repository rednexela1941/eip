package adapter

import (
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/network"
)

type RequestContext struct {
	*network.Info
	PortInstanceID cip.UINT
}

func NewRequestContext(info *network.Info, portInstanceID cip.UINT) *RequestContext {
	return &RequestContext{
		Info:           info,
		PortInstanceID: portInstanceID,
	}
}
