// Package adapter implements the core interfaces and helpers required to create an EtherNet/IP compatible adapter/slave device.
package adapter

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rednexela1941/eip/pkg/encap"
	"github.com/rednexela1941/eip/pkg/identity"
)

const (
	DefaultNetworkTickInterval = 10 * time.Millisecond
)

type (
	Adapter struct {
		*identity.Identity
		Classes           ClassMap
		Logger            *log.Logger
		Connections       *ConnectionStore
		ConnectionPoints  []ConnectionPoint
		AssemblyInstances []*AssemblyInstance

		sessionHandleOffset encap.SessionHandle
		NetworkTickInterval time.Duration

		listenerParams []listenParams
	}
	_Adapter = Adapter
)

func New(ident *identity.Identity) *Adapter {
	const logFlags = log.Lshortfile | log.Lmsgprefix | log.Ltime | log.Lmicroseconds
	a := new(_Adapter)
	a.Logger = log.New(os.Stderr, "", logFlags)
	a.Identity = ident
	a.Classes = make(ClassMap)
	a.Connections = NewConnectionStore(32)
	a.ConnectionPoints = make([]ConnectionPoint, 0, 2)
	a.listenerParams = make([]listenParams, 0, 1)
	a.NetworkTickInterval = DefaultNetworkTickInterval

	a.InitDefaultIdentityObject()
	a.InitDefaultMessageRouterObject()
	a.InitDefaultConnectionManagerObject()
	return a
}

func (self *Adapter) SetNetworkTickInterval(duration time.Duration) error {
	if duration < 0 {
		return fmt.Errorf("invalid duration %s", duration.String())
	}
	self.NetworkTickInterval = duration
	return nil
}

func (self *_Adapter) Handle(c *RequestContext, p encap.Request) (encap.Reply, error) {
	// self.Logger.Printf("-> %s", p.GetCommand().String()) // TODO: also include network info.
	reply, err := self._Handle(c, p)
	self.Logger.Printf("%s on %s\n", p.GetCommand().String(), c.String())
	return reply, err
}

func (self *_Adapter) _Handle(c *RequestContext, p encap.Request) (encap.Reply, error) {
	switch p.GetCommand() {
	case encap.NOP:
		// Vol 2: 2-4.1
		return self._HandleNOP(c, p.(encap.NOPRequest))
	case encap.ListIdentity:
		// Vol 2: 2-4.2
		return self._HandleListIdentity(c, p.(encap.ListIdentityRequest))
	case encap.ListInterfaces:
		// Vol 2: 2-4.3
		return self._HandleListInterfaces(c, p.(encap.ListInterfacesRequest))
		// return self.HandleListInterfaces(c, p)
	case encap.ListServices:
		// Vol 2: 2-4.6
		return self._HandleListServices(c, p.(encap.ListServicesRequest))
	case encap.RegisterSession:
		// Vol 2: 2-4.4
		return self._HandleRegisterSession(c, p.(encap.RegisterSessionRequest))
	case encap.UnRegisterSession:
		// Vol 2: 2-4.5
		return self._HandleUnregisterSession(c, p.(encap.UnregisterSessionRequest))
	case encap.SendRRData:
		// Vol 2: 2-4.7
		return self._HandleSendRRData(c, p.(encap.SendRRDataRequest))
	case encap.SendUnitData:
		// Vol 2: 2-4.8
		return self._HandleSendUnitData(c, p.(encap.SendUnitDataRequest))
	default:
		return nil, fmt.Errorf("unknown encapsulation command: %s", p.GetCommand().String())
	}
}
