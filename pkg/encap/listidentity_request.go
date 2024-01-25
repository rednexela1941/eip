package encap

import (
	"math/rand"
	"time"

	"github.com/rednexela1941/eip/pkg/cip"
)

/*
A connection originator may use the ListIdentity command to locate and identify potential
targets.  This command shall be sent as a unicast message using TCP or UDP, or as a broadcast
message using UDP and does not require that a session be established. The reply shall always
be sent as a unicast message.
When received as a broadcast message, the receiving device shall delay for a pseudo-random
period of time prior to sending the reply as specified in section 2-4.2.3.  Delaying before
sending the reply helps to spread out any resulting ARP requests and ListIdentity replies from
target devices on the network.
*/
// Volume 2: Table 2-4.2
type ListIdentityRequest interface {
	Request
	// See Volume 2: 2-4.3.3 (MaxResponseDelay)
	GetMaxResponseDelay() cip.UINT
	GetResponseDelay() time.Duration
}

type _ListIdentityRequest Packet

func (self *Packet) ToListIdentityRequest() (ListIdentityRequest, error) {
	return (*_ListIdentityRequest)(self), nil
}

func (self *_ListIdentityRequest) GetMaxResponseDelay() cip.UINT {
	senderContext := self.SenderContext

	maxResponseDelay := cip.UINT(0)
	maxResponseDelay |= cip.UINT(senderContext[0])
	maxResponseDelay |= cip.UINT(senderContext[1]) << 8

	if maxResponseDelay == 0 {
		maxResponseDelay = 2000
	}
	if maxResponseDelay < 500 {
		maxResponseDelay = 500
	}
	return maxResponseDelay
}

func (self *_ListIdentityRequest) GetResponseDelay() time.Duration {
	maxDelay := int(self.GetMaxResponseDelay()) // milliseconds
	rando := cip.UINT(rand.Intn(int(maxDelay)))
	return time.Millisecond * time.Duration(rando)
}
