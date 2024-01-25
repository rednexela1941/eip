package adapter

import (
	"fmt"
	"net"

	"github.com/rednexela1941/eip/pkg/cpf"
)

// ChannelSockaddrs represents the sockaddrs
// for sending and receiving IO data.
type ChannelSockaddrs struct {
	OtoT cpf.SockaddrInfo
	TtoO cpf.SockaddrInfo
}

func (self *ChannelSockaddrs) GetTtoOUDPAddr() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   self.TtoO.SinAddr.ToIP(),
		Port: int(self.TtoO.SinPort),
	}
}

func (self *ChannelSockaddrs) String() string {
	return fmt.Sprintf("%s, %s", self.OtoT.String(), self.TtoO.String())
}
