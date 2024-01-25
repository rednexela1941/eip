package adapter

import (
	"fmt"
	"math/rand"

	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cm"
	"github.com/rednexela1941/eip/pkg/encap"
)

// ConnectionStore is used in the adapter to store active connections
type ConnectionStore struct {
	Ring          []*Connection
	IncarnationID uint32 // see Volume 2: 3-3.7.1.2
	idOffset      uint16
}

func NewConnectionStore(maxNumberOfConnections int) *ConnectionStore {
	return &ConnectionStore{
		Ring:          make([]*Connection, maxNumberOfConnections),
		IncarnationID: rand.Uint32() << 16,
	}
}

func (self *ConnectionStore) NewConnectionID() cip.UDINT {
	self.idOffset++
	return self.IncarnationID + uint32(self.idOffset)
}

func (self *ConnectionStore) AddConnection(conn *Connection) error {
	for i, c := range self.Ring {
		if c == nil {
			self.Ring[i] = conn
			return nil
		}
	}
	return fmt.Errorf("too many connections (%d)", len(self.Ring))
}

func (self *ConnectionStore) HasNonListenOnlyConnectionTo(
	instance *AssemblyInstance,
) bool {
	for _, c := range self.Ring {
		if c == nil {
			continue
		}
		cp := c.Point
		if cp == nil || cp.Type == ListenOnly {
			continue
		}
		if cp.Input == instance {
			return true
		}
	}
	return false
}

func (self *ConnectionStore) GetConnectionByOtoTID(id cip.UDINT) (conn *Connection, ok bool) {
	for _, c := range self.Ring {
		if c == nil {
			continue
		}
		if c.OtoTNetworkConnectionID == id {
			return c, true
		}
	}
	return nil, false
}

func (self *ConnectionStore) ConnectionExists(triad cm.Triad) bool {
	_, ok := self.GetConnection(triad)
	return ok
}

func (self *ConnectionStore) GetConnection(triad cm.Triad) (conn *Connection, ok bool) {
	for _, c := range self.Ring {
		if c == nil {
			continue
		}
		if c.Triad == triad {
			return c, true
		}
	}
	return nil, false
}

// RemoveConnectionBySessionHandle: for cleanup of class 2 and class 3
// connections after timeout.
func (self *ConnectionStore) RemoveConnectionBySessionHandle(sh encap.SessionHandle) bool {
	if sh == 0 {
		return false
	}
	for i, c := range self.Ring {
		if c == nil {
			continue
		}
		if c.SessionHandle == sh {
			self.Ring[i] = nil
			return true
		}
	}
	return false
}

func (self *ConnectionStore) RemoveConnection(triad cm.Triad) bool {
	for i, c := range self.Ring {
		if c == nil || c.Triad != triad {
			continue
		}
		// if err := c.CloseTCP(); err != nil {
		// 	log.Println(err)
		// }
		self.Ring[i] = nil
		return true
	}
	return false // connection not found.
}
