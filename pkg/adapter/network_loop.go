package adapter

import (
	"time"

	"github.com/rednexela1941/eip/pkg/cm"
	"github.com/rednexela1941/eip/pkg/cpf"
)

// StartNetworkLoop
// blocking function to run the EtherNet/IP server.
func (self *Adapter) StartNetworkLoop() error {
	for _, p := range self.listenerParams {
		if err := self.ListenOn(p.addr, p.iface, p.pctx); err != nil {
			return err
		}
	}
	tick := self.NetworkTickInterval

	for {
		if err := self._CloseStaleConnections(); err != nil {
			self.Logger.Println(err)
		}
		if err := self._SendIOData(tick); err != nil {
			self.Logger.Println(err)
		}
		time.Sleep(tick)
	}
}

func (self *Adapter) _SendIOData(tick time.Duration) error {
	for _, c := range self.Connections.Ring {
		if c == nil || !c.IsIO() {
			continue
		}
		// only class 0 and class 1 connections here.
		// now, check if we should send something.
		if !c.IsTimeToSend(tick) {
			continue
		}

		if c.TransportClassAndTrigger.ProductionTrigger() != cm.Cyclic {
			// See Volume 1: Table 3-4.14
			self.Logger.Fatal("todo, other trigger actions")
		}
		// check triggers in time to send.

		channels := c.Channels
		if channels == nil {
			self.Logger.Fatal("channels is nil for Class0/Class1 connection")
		}

		cp := c.Point
		if cp == nil {
			self.Logger.Fatal("connection point is nil on Class0/Class1 connection")
		}

		input := cp.Input
		if input == nil {
			continue // nothing to send.
		}

		// now, what format is the data in?
		writer := cpf.NewWriter()

		c.TtoOSequenceNumber++
		c.TtoOEncapsulationSequenceNumber++

		// does this apply to both class 0 and class 1 packets?
		writer.AddSequencedAddressItem(
			c.TtoONetworkConnectionID,
			c.TtoOEncapsulationSequenceNumber,
		)
		dataItem := writer.AddIODataItem()

		ioConn := c.Info.IOConn
		if ioConn == nil {
			self.Logger.Fatal("io conn is nil")
		}

		c.writeTtoOIODataHeader(dataItem)

		if err := input.WriteTo(dataItem); err != nil {
			self.Logger.Println(err)
			continue
		}
		if dataItem.Error() != nil {
			self.Logger.Fatal(dataItem.Error())
			continue
		}

		data, err := writer.Encode()
		if err != nil {
			self.Logger.Fatal(err)
			continue
		}

		destAddr := channels.GetTtoOUDPAddr()

		if _, err := ioConn.WriteTo(data, nil, destAddr); err != nil {
			self.Logger.Println(err)
			continue
		}
		c.UpdateTtoOTimestamp()
	}
	return nil
}

func (self *Adapter) _CloseStaleConnections() error {
	for i, c := range self.Connections.Ring {
		if c == nil || !c.IsOtoTTimedOut() {
			continue
		}
		self.Logger.Printf("%s timed out", c.String())
		if err := c.CloseTCP(); err != nil {
			self.Logger.Println(err)
		}
		self.Connections.Ring[i] = nil // clear out connection.
	}
	return nil
}

// CloseAllConnections: cleanup all connections
// will be used inside of reset function to speed things up.
// in conformance tests.
func (self *Adapter) CloseAllConnections() error {
	for i, c := range self.Connections.Ring {
		if c == nil {
			continue
		}
		if c.Conn != nil {
			if err := c.Conn.Close(); err != nil {
				self.Logger.Println(err)
			}
			self.Connections.Ring[i] = nil
			continue
		}
		// TODO: make a nicer close connection function.
	}
	return nil
}
