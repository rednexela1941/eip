package cm

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/epath"
)

type ForwardCloseRequestHeader struct {
	PriorityTimeTick cip.BYTE  // see 3-5.6.1.2.1
	TimeoutTicks     cip.USINT // see 3-5.6.1.2.1
	Triad
	ConnectionPathSizeWords cip.USINT
	Reserved                cip.USINT
}

// Volume 1: Table 3-5.24
type ForwardCloseRequest struct {
	ForwardCloseRequestHeader
	ConnectionPath epath.PaddedEPATH
}

func NewForwardCloseRequest(r bbuf.Reader) *ForwardCloseRequest {
	f := new(ForwardCloseRequest)
	r.Rl(&f.ForwardCloseRequestHeader)
	if f.ConnectionPathSizeWords > 0 {
		data := make([]byte, int(f.ConnectionPathSizeWords)*2)
		r.Rl(&data)
	}
	return f
}
