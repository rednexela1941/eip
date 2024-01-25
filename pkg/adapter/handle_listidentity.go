package adapter

import (
	"time"

	"github.com/rednexela1941/eip/pkg/cpf"
	"github.com/rednexela1941/eip/pkg/encap"
	"github.com/rednexela1941/eip/pkg/network"
)

// Vol 2: 2-4.2
func (self *_Adapter) _HandleListIdentity(c *RequestContext, p encap.ListIdentityRequest) (encap.ListIdentityReply, error) {
	if p.GetOptions() != 0 {
		// discard non-zero options requests.
		return nil, nil
	}

	if c.IsMulticast() {
		// random delay
		time.Sleep(p.GetResponseDelay())
	}

	reply := encap.NewListIdentityReply()
	reply.SetHeader(p.GetHeader())

	// see Volume 2: Table 2-4.4
	reply.AddIdentityItem(
		&cpf.SockaddrInfo{
			SinFamily: 2,               // AF_INET
			SinPort:   network.TCPPort, // 0xAF12
			SinAddr:   c.Local.IP,
		},
		self.Identity,
		encap.ProtocolVersion,
	)

	return reply, nil
}
