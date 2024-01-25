package adapter

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cm"
	"github.com/rednexela1941/eip/pkg/cpf"
	"github.com/rednexela1941/eip/pkg/encap"
)

const ConnectionManagerObjectRevision cip.UINT = 1

func (self *_Adapter) InitDefaultConnectionManagerObject() {
	c := self.AddClass("Connection Manager", cip.ConnectionManagerClassCode, ConnectionManagerObjectRevision)

	c.AddAttribute(1, "ClassRevision", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(ConnectionManagerObjectRevision) },
	)
	c.AddAttribute(2, "MaxInstance", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestInstanceID()) },
	)
	c.AddAttribute(3, "NumInstances", cip.UINTSize).OnGet(
		GetSingle,
		func(res Response) { res.Wl(c.NumberOfInstances()) },
	)
	c.AddAttribute(6, "MaxClassAttributeID", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestAttributeID()) },
	)
	c.AddAttribute(7, "MaxInstanceAttributeID", cip.UDINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestInstanceAttributeID()) },
	)
	c.addDefaultGetAttributesAll()

	i := c.AddInstance(1)

	// may not need this anymore.
	// might be required for the correct
	// attribute error.
	// see if we should generalize this out to the rest of
	// the CallService functions.
	nullService := func(req *Request, res Response) {
		// for whatever reason, these services are supported, but without attributes.
		res.SetGeneralStatus(cip.StatusAttributeNotSupported)
	}

	i.OnService(cip.GetAttributeSingle, nullService)
	i.OnService(cip.ForwardOpen, self.HandleForwardOpen)
	i.OnService(cip.LargeForwardOpen, self.HandleLargeForwardOpen)
	i.OnService(cip.ForwardClose, self.HandleForwardClose)
}

func (self *Adapter) HandleForwardClose(req *Request, res Response) {
	// By Volume 2: 3-3.10
	// We must also validate that the forward close request
	// comes from the same IP address
	// as the person who opened the connection.
	// also, should close connections based on UnregisterSession.

	r := bbuf.New(req.RequestData)
	freq := cm.NewForwardCloseRequest(r)
	if r.Error() != nil {
		res.SetGeneralStatus(cip.StatusNotEnoughData)
		return
	}
	if r.Len() > 0 {
		// too much data.
		res.SetGeneralStatus(cip.StatusTooMuchData)
		return
	}

	conn, ok := self.Connections.GetConnection(freq.Triad)
	if !ok {
		self.Logger.Println("connection doesn't exist")
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(0x107) // Target connection not found.
		return
	}

	if conn.Info.Remote.IP != req.Info.Remote.IP {
		self.Logger.Println("IP address mismatch")
		res.SetGeneralStatus(cip.StatusPriveledgeViolation)
		return
	}

	if !self.Connections.RemoveConnection(freq.Triad) {
		self.Logger.Fatal("failed to remove connection")
	}

	fres := new(cm.ForwardCloseResponse)
	fres.Triad = freq.Triad
	res.Wl(&fres.Triad)
	res.Wl(fres.ApplicationReplySizeWords)
	res.Wl(fres.Reserved)
}

// See Volume 2: Table 3-3.3 for usage of SockaddrInfo items in forward open.
func (self *Adapter) HandleLargeForwardOpen(req *Request, res Response) {
	self._HandleForwardOpen(req, res)
}

// See Volume 2: Table 3-3.3 for usage of SockaddrInfo items in forward open.
// there is another table that we should look at -- just need to find it.
func (self *Adapter) HandleForwardOpen(req *Request, res Response) {
	self._HandleForwardOpen(req, res)
}

// See Volume 1: 3-5.6.3 "Forward_Open and Large_Forward_Open requests"
// See Volume 1: 3-6.1 For Realtime Format information.
// See Volume 2: Table 3-3.3 for usage of SockaddrInfo items in forward open.
// See Volume 1: 3-5.6.1.1.1 for connection timeout handling.
// See Volume 1: 3-5.6.1 (Chapter) for service parameter descriptions.
// See Volume 1: Table 3-5.16 for connection path ordering.
// See Volume 2: 3-3.9 Forward_Open for CIP Transport Class 0 and Class 1 Connections.
// See Volume 1: 3-6 "Application Connection type using Class 0 and Class 1"
func (self *_Adapter) _HandleForwardOpen(req *Request, res Response) {
	var freq *cm.SharedForwardOpenRequest
	r := bbuf.New(req.RequestData)

	switch req.Service {
	case cip.ForwardOpen:
		freq = cm.NewForwardOpenRequest(r)
	case cip.LargeForwardOpen:
		freq = cm.NewLargeForwardOpenRequest(r)
	default:
		self.Logger.Fatalf("invalid %s", req.Service.String())
	}

	if r.Error() != nil {
		// not enough data.
		self.Logger.Println(r.Error())
		res.SetGeneralStatus(cip.StatusNotEnoughData)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return
	}

	if r.Len() > 0 {
		// too much data
		self.Logger.Printf("%d bytes remaining", r.Len())
		res.SetGeneralStatus(cip.StatusTooMuchData)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return
	}

	if self.Connections.ConnectionExists(freq.Triad) {
		// already have a connection.
		// See Volume 1: 3-5.6.3
		// for more details.
		// has to be a null forward open.

		if !freq.IsNull() {
			self.Logger.Println("dont know how to deal with matching null")
			res.SetGeneralStatus(cip.StatusCommunicationProblem)
			res.AddAdditionalStatus(cip.DuplicateForwardOpen)
			cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
			return
		}
		self.Logger.Println("matching, but null -- not implemented.")
		return // don't know what to do here.
	}

	// non-matching
	if freq.IsNull() {
		// See Volume 1: 3-5.6.2.2.1
		self.Logger.Println("don't know how to deal with null")
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.NullForwardOpenNotSupported)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return
	}

	if freq.OtoTParameters.IsReserved() || freq.TtoOParameters.IsReserved() {
		self.Logger.Println("reserved connection type")
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		if freq.OtoTParameters.IsReserved() {
			res.AddAdditionalStatus(cip.InvalidOtoTConnectionType)
		}
		if freq.TtoOParameters.IsReserved() {
			res.AddAdditionalStatus(cip.InvalidTtoOConnectionType)
		}
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return
	}

	if !checkTransportClassAndTrigger(freq, res) {
		self.Logger.Printf(
			"invalid TransportClassAndTrigger: %d\n",
			freq.TransportClassAndTrigger,
		)
		return
	}

	if !self.checkConnectionPathAndKey(freq, res) {
		return
	}

	// Connection ID Assignment.
	// See Volume 2: Table 3-3.2 Network Connection ID Selection
	self.assignConnectionIDs(freq)

	switch freq.TransportClassAndTrigger.TransportClass() {
	case cm.Class0, cm.Class1:
		self._HandleClass0And1ForwardOpen(req, freq, res)
		return
	case cm.Class2, cm.Class3:
		self._HandleClass2And3ForwardOpen(req, freq, res)
		return
	default:
		// shouldn't get here, transport class and trigger is checked above.
		self.Logger.Fatalln("fatal: invalid TransportClassAndTrigger")
	}
}

func (self *_Adapter) _HandleClass0And1ForwardOpen(
	req *Request,
	freq *cm.SharedForwardOpenRequest,
	res Response,
) {
	cp, err := self.GetMatchingConnectionPoint(freq)

	if err != nil || cp == nil {
		self.Logger.Println(err, "connection point not found")
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.InvalidConnectionPath)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return
	}

	transportClass := freq.TransportClassAndTrigger

	// check sizes.
	if cp.InputSize(transportClass) != freq.TtoOParameters.Size {
		self.Logger.Printf("TtoO size mismatch: got=%d need=%d\n", freq.TtoOParameters.Size, cp.InputSize(transportClass))
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.InvalidTtoONetworkConnectionSize)
		res.AddAdditionalStatusNoReplace(cp.InputSize(transportClass)) // attach correct status.
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return
	}

	if cp.OutputSize(transportClass) != freq.OtoTParameters.Size {
		self.Logger.Printf("OtoT size mismatch: got=%d, need=%d\n", freq.OtoTParameters.Size, cp.OutputSize(transportClass))
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.InvalidOtoTNetworkConnectionSize)
		res.AddAdditionalStatusNoReplace(cp.OutputSize(transportClass)) // attach correct size.
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return
	}

	channels, err := getChannelSockaddrs(req, freq)
	if err != nil {
		self.Logger.Println(err)
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.ParameterErrorUnconnectedService)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return
	}

	if cp.Type == ListenOnly && !self.Connections.HasNonListenOnlyConnectionTo(cp.Input) {
		// TODO: cleanup on close connection.
		self.Logger.Println("no connections to attach listeonly")
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.NonListenOnlyConnectionNotOpened)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return
	}

	if !self.CheckRPI(freq, res) {
		return
	}

	writer, err := encap.GetItemWriter(res.Parent())
	if err != nil {
		self.Logger.Fatal(err)
		return
	}

	// Make sure that OtoT comes first. not required, but some devices assume this.
	// Apparently we have to flip them around?
	// so OtoT gets send as TtoO
	// TODO: check this.
	writer.AddSockaddrInfoItem(
		cpf.SockaddrInfoOtoT,
		channels.OtoT,
	)
	writer.AddSockaddrInfoItem(
		cpf.SockaddrInfoTtoO,
		channels.TtoO,
	)

	conn := NewClass0Or1Connection(req.Info, freq, cp, channels)

	if err := self.tryToAddConnection(conn, res); err != nil {
		self.Logger.Println(err)
		return
	}

	sendForwardOpenReplySuccess(freq, res)
}

func (self *_Adapter) _HandleClass2And3ForwardOpen(
	req *Request,
	freq *cm.SharedForwardOpenRequest,
	res Response,
) {
	firstPath := freq.ConnectionPath.ApplicationPaths[0]
	if firstPath.ClassID != cip.MessageRouterClassCode || firstPath.InstanceID != 1 {
		// figure out what error should go here.
		self.Logger.Println("invalid class2/class3 path")
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.TargetOutOfConnections)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return

	}

	// Add Class2/Class3 Connection.
	conn := NewClass2Or3Connection(
		req.Info,
		freq,
		res.Parent().PeekSessionHandle(),
	)

	if err := self.tryToAddConnection(conn, res); err != nil {
		self.Logger.Println(err)
		return
	}

	sendForwardOpenReplySuccess(freq, res)
}

// write successful forward open response to res.
func sendForwardOpenReplySuccess(
	freq *cm.SharedForwardOpenRequest,
	res Response,
) {
	fres := new(cm.ForwardOpenResponse)
	fres.Triad = freq.Triad

	fres.OtoTNetworkConnectionID = freq.OtoTNetworkConnectionID
	fres.TtoONetworkConnectionID = freq.TtoONetworkConnectionID

	fres.OtoTAPI = freq.OtoTParameters.RPI
	fres.TtoOAPI = freq.TtoOParameters.RPI

	res.Wl(&fres.ForwardOpenResponseHeader)
}
