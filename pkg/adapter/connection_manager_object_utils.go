package adapter

import (
	"fmt"
	"time"

	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cm"
	"github.com/rednexela1941/eip/pkg/cpf"
	"github.com/rednexela1941/eip/pkg/encap"
	"github.com/rednexela1941/eip/pkg/epath"
	"github.com/rednexela1941/eip/pkg/network"
)

func (self *_Adapter) _ValidRPI(t time.Duration) bool {
	return t > 0 && t >= self.networkTickInterval
}

func (self *_Adapter) CheckRPI(freq *cm.SharedForwardOpenRequest, res Response) bool {
	if !self._ValidRPI(freq.OtoTParameters.GetRPI()) || !self._ValidRPI(freq.TtoOParameters.GetRPI()) {
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.RPINotAccepted)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return false
	}
	return true
}

// See Volume 2: Table 3-3.2 "Network Connection ID Selection"
func (self *_Adapter) assignConnectionIDs(freq *cm.SharedForwardOpenRequest) {

	o2t := freq.OtoTParameters.Type
	t2o := freq.TtoOParameters.Type

	if o2t.IsPointToPoint() {
		// freq.OtoTNetworkConnectionID =
		freq.OtoTNetworkConnectionID = self.Connections.NewConnectionID()
	}

	if t2o.IsMulticast() {
		freq.TtoONetworkConnectionID = self.Connections.NewConnectionID()
	}
}

func getChannelSockaddrs(
	req *Request,
	freq *cm.SharedForwardOpenRequest,
) (*ChannelSockaddrs, error) {
	// okay, now we need some stuff.
	// See Volume 2: Table 3-3.3 "Sockaddr Info Usage"
	sockets := new(ChannelSockaddrs)

	o2t := freq.OtoTParameters.Type
	t2o := freq.TtoOParameters.Type

	o2tItem, _ := encap.GetSockaddrInfoReader(req.Parent, cpf.SockaddrInfoOtoT)
	t2oItem, _ := encap.GetSockaddrInfoReader(req.Parent, cpf.SockaddrInfoTtoO)

	// TODO: check that these are valid.
	// must have AF_INET = 2
	if o2t.IsPointToPoint() && t2o.IsMulticast() {
		// ignore request values.
		sockets.OtoT = cpf.NewSockaddrInfo(
			req.Info.Local.IP,
			network.UDPIOPort,
		)
		sockets.TtoO = cpf.NewSockaddrInfo(
			req.Info.GetMulticastAddress(),
			network.UDPIOPort,
		)
		return sockets, nil
	}
	if o2t.IsMulticast() && t2o.IsPointToPoint() {
		if o2tItem == nil {
			return nil, fmt.Errorf("missing required OtoT sockaddr info")
		}
		// todo, check valid.
		sockets.OtoT = o2tItem.GetSockaddrInfo()
		if !sockets.OtoT.SinAddr.IsMulticast() {
			return nil, fmt.Errorf("sockaddr is not multicast")
		}
		if t2oItem != nil {
			sockets.TtoO = t2oItem.GetSockaddrInfo()
			// force using same IP.
			sockets.TtoO.SinAddr = req.Info.Remote.IP
		} else {
			sockets.TtoO = cpf.NewSockaddrInfo(
				req.Info.Remote.IP,
				network.UDPIOPort,
			)
		}
		return sockets, nil
	}
	if o2t.IsPointToPoint() && t2o.IsPointToPoint() {
		sockets.OtoT = cpf.NewSockaddrInfo(
			req.Info.Local.IP,
			network.UDPIOPort,
		)
		if t2oItem != nil {
			sockets.TtoO = t2oItem.GetSockaddrInfo()
			// force using same IP.
			sockets.TtoO.SinAddr = req.Info.Remote.IP
		} else {
			sockets.TtoO = cpf.NewSockaddrInfo(
				req.Info.Remote.IP,
				network.UDPIOPort,
			)
		}
		return sockets, nil
	}
	if o2t.IsMulticast() && t2o.IsMulticast() {
		if o2tItem == nil {
			return nil, fmt.Errorf("missing required OtoT sockaddr info")
		}
		sockets.OtoT = o2tItem.GetSockaddrInfo()
		// TODO: validate multicast address
		if !sockets.OtoT.SinAddr.IsMulticast() {
			return nil, fmt.Errorf("sockaddr is not multicast")
		}
		sockets.TtoO = cpf.NewSockaddrInfo(
			req.Info.GetMulticastAddress(),
			network.UDPIOPort,
		)
		return sockets, nil
	}
	return nil, fmt.Errorf("unknown conn types (%d, %d)", o2t, t2o)
}

// this function will write appropriate error message.
func (self *_Adapter) tryToAddConnection(conn *Connection, res Response) error {
	if err := self.Connections.AddConnection(conn); err != nil {
		self.Logger.Println(err)
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.TargetOutOfConnections)
		return err
	}
	self.Logger.Printf("Added New %s\n", conn.String())
	return nil
}

// CheckConnectionPath: return true if valid.
// See Volume 1: Table 3-5.16 for application paths.
func (self *_Adapter) CheckConnectionPaths(freq *cm.SharedForwardOpenRequest, res Response) bool {
	aps := freq.ConnectionPath.ApplicationPaths

	writeErr := func() {
		self.Logger.Println("invalid connection path", freq.ConnectionPath.String())
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.InvalidConnectionPath)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
	}

	if len(aps) == 0 {
		writeErr()
		return false
	}

	for _, a := range aps {
		if !self.checkObjectExists(&a) {
			writeErr()
			return false
		}
	}

	// TODO: check that it somehow matches a connection point.
	return true
}

func (self *_Adapter) checkConnectionPathAndKey(
	freq *cm.SharedForwardOpenRequest,
	res Response,
) (ok bool) {
	_, err := freq.ConnectionPath.Parse()
	if err != nil {
		self.Logger.Println(err)
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		res.AddAdditionalStatus(cip.InvalidConnectionPath)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
		return false
	}
	// validate electronic key.
	ekey := freq.ConnectionPath.ElectronicKey
	if ekey != nil {
		if ekey.Format == epath.Format5 {
			// no Format5, it seems (from CT20)
			res.SetGeneralStatus(cip.StatusCommunicationProblem)
			res.AddAdditionalStatus(cip.InvalidConnectionPath)
			cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
			return false
		}

		valid, addStatus := ValidateElectronicKey(self.Identity, ekey)
		if !valid {
			res.SetGeneralStatus(cip.StatusCommunicationProblem)
			for _, s := range addStatus {
				res.AddAdditionalStatus(s)
			}
			cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
			return false
		}
	}

	return self.CheckConnectionPaths(freq, res)
}

func checkTransportClassAndTrigger(freq *cm.SharedForwardOpenRequest, res Response) (ok bool) {
	tct := freq.TransportClassAndTrigger
	ok = true
	if !tct.TransportClassValid() {
		ok = false
		res.AddAdditionalStatus(cip.TransportClassNotSupported)
	}
	if !tct.ProductionTriggerValid() {
		ok = false
		res.AddAdditionalStatus(cip.TtoOTriggerNotSupported)
	}
	if !ok {
		res.SetGeneralStatus(cip.StatusCommunicationProblem)
		cm.WriteUnsuccessfulForwardOpenResponse(res, freq)
	}
	return
}

func (self *_Adapter) checkObjectExists(a *epath.ApplicationPath) bool {
	if a.ClassID == cip.AssemblyClassCode {
		if a.InstanceID == 0 && a.ConnectionPoint != 0 {
			a.InstanceID = a.ConnectionPoint
		}
	}

	c, ok := self.Classes[a.ClassID]
	if !ok {
		return false
	}

	if a.InstanceID != 0 {
		_, ok = c.Instances[a.InstanceID]
		if !ok {
			return false
		}
	}

	return true
}
