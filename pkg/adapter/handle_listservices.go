package adapter

import (
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cpf"
	"github.com/rednexela1941/eip/pkg/encap"
)

type CapabilityFlags = cip.UINT

const (
	CapabilityFlagEtherNetIP            CapabilityFlags = 1 << 4
	CapabilityFlagClass0And1Connections CapabilityFlags = 1 << 7
)

// Volume 2: 2-4.6.3
func (self *_Adapter) _HandleListServices(c *RequestContext, p encap.ListServicesRequest) (encap.ListServicesReply, error) {
	if p.GetOptions() != 0 {
		// discard non-zero options requests.
		return nil, nil
	}

	reply := encap.NewListServicesReply()
	reply.SetHeader(p.GetHeader())

	item := reply.AddItem()
	item.SetTypeID(cpf.ListServicesResponse)
	item.Wl(encap.ProtocolVersion)
	item.Wl(CapabilityFlagEtherNetIP | CapabilityFlagClass0And1Connections)

	nameOfService := make([]byte, 16)
	copy(nameOfService, []byte("Communications"))
	if _, err := item.Write(nameOfService); err != nil {
		self.Logger.Println(err)
	}

	return reply, nil
}
