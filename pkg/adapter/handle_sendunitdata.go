package adapter

import (
	"github.com/rednexela1941/eip/pkg/cpf"
	"github.com/rednexela1941/eip/pkg/encap"
)

// The SendUnitData command shall send encapsulated connected messages.  A reply shall not be
// returned.
// NOTE: When used to encapsulate the CIP, the SendUnitData command is used to send CIP
// connected data in both the O to T and T to O directions.
// Volume 2: 2-4.8
// See Volume 2: Table 2-4.17
func (self *_Adapter) _HandleSendUnitData(c *RequestContext, request encap.SendUnitDataRequest) (encap.SendUnitDataReply, error) {
	if request.GetOptions() != 0 {
		// discard non-zero options requests.
		return nil, nil
	}

	reply := encap.NewSendUnitDataReply()
	reply.SetHeader(request.GetHeader())

	sessionHandle := request.GetSessionHandle()

	if (c.IsUDP() && sessionHandle != 0) || !self.IsValidSession(c, sessionHandle) {
		reply.SetStatus(encap.StatusInvalidSessionHandle)
		return reply, nil
	}

	rp := request.GetEncapsulatedPacket()

	if rp.GetItemCount() < 2 {
		// not enogh items.
		reply.SetStatus(encap.StatusInvalidLength)
		return reply, nil
	}

	addrItem, ok := rp.GetItem(0).(cpf.ConnectedAddressItemReader)
	if !ok {
		reply.SetStatus(encap.StatusInvalidCommand)
		return reply, nil
	}

	dataItem, ok := rp.GetItem(1).(cpf.ConnectedDataItemReader)
	if !ok {
		reply.SetStatus(encap.StatusInvalidCommand)
		return reply, nil
	}

	addrItem.GetTypeID() // silence compiler.
	dataItem.GetTypeID() // silence compiler.

	// TODO: reply doesn't need connection identifier?
	// Also, go and make sure the connection with ID exists.
	// Then we can set the connection ID to something appropriate.
	connection, ok := self.Connections.GetConnectionByOtoTID(
		addrItem.GetConnectionID(),
	)

	if !ok {
		// TODO: what would be the correct error code here?
		reply.SetStatus(encap.StatusInvalidCommand)
		return reply, nil
	}

	connection.UpdateOtoTTimestamp()
	// reply.AddConnectedAddressItem(addrItem.GetConnectionIdentifier())
	// find connection and get data.
	reply.AddConnectedAddressItem(connection.TtoONetworkConnectionID)

	// TODO: check sequence number and continue.
	seqNumber := dataItem.GetSequenceCount()
	if seqNumber == connection.OtoTSequenceNumber {
		// resend last message
		if connection.LastSentPacket != nil {
			return connection.LastSentPacket, nil
		}
		self.Logger.Println("sequence mismatch with no cached reply", seqNumber)
	}

	connection.OtoTSequenceNumber = seqNumber
	connection.TtoOSequenceNumber = seqNumber

	resDataItem := reply.AddConnectedDataItem(seqNumber)

	self.Route(
		NewRequest(request, c, dataItem.GetMessageRouterRequest()),
		NewResponse(reply, resDataItem),
	)

	connection.UpdateTtoOTimestamp()

	connection.LastSentPacket = reply

	return reply, nil
}
