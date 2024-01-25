package adapter

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/encap"
)

// Volume 2: 2-4.3
func (self *_Adapter) _HandleListInterfaces(c *RequestContext, p encap.ListInterfacesRequest) (encap.ListInterfacesReply, error) {
	if p.GetOptions() != 0 {
		// discard packets with non-zero options.
		return nil, nil
	}

	if p.GetLength() != 0 {
		return nil, fmt.Errorf("invalid length (%d) for %s", p.GetLength(), p.GetCommand().String())
	}

	reply := encap.NewListInterfacesReply()
	reply.SetHeader(p.GetHeader())

	reply.Wl(cip.UINT(0)) // zero item reply.

	return reply, nil
}
