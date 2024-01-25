package adapter

import (
	"github.com/rednexela1941/eip/pkg/cpf"
	"github.com/rednexela1941/eip/pkg/encap"
)

// Volume 2: 2-4.7
// NOTE: When used to encapsulate the CIP, the SendRRData request and response are used to
// send encapsulated UCMM messages (unconnected messages).  See chapter 3 for more detail.
// Volume 2: 3-2.1 (UCMM messages)
// Reply See Volume 2: Table 3-2.2
func (self *_Adapter) _HandleSendRRData(c *RequestContext, request encap.SendRRDataRequest) (encap.SendRRDataReply, error) {
	if request.GetOptions() != 0 {
		// discard non-zero options requests.
		return nil, nil
	}

	reply := encap.NewSendRRDataReply()
	reply.SetHeader(request.GetHeader())

	rp := request.GetEncapsulatedPacket()

	if rp.GetItemCount() < 2 {
		// not enogh items.
		reply.SetStatus(encap.StatusInvalidLength)
		return reply, nil
	}

	addrItem, ok := rp.GetItem(0).(cpf.NullAddressItemReader)
	if !ok {
		reply.SetStatus(encap.StatusInvalidCommand)
		return reply, nil
	}

	dataItem, ok := rp.GetItem(1).(cpf.UnconnectedDataItemReader)
	if !ok {
		reply.SetStatus(encap.StatusInvalidCommand)
		return reply, nil
	}

	noop(addrItem) // silence compiler
	noop(dataItem) // silence compiler

	reply.AddNullAddressItem()

	resDataItem := reply.AddUnconnectedDataItem()

	self.Route(
		NewRequest(request, c, dataItem.GetMessageRouterRequest()),
		NewResponse(reply, resDataItem),
	)

	// TODO: route
	return reply, nil
}
