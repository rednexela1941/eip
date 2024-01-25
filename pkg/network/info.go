package network

import (
	"fmt"
	"net"

	"golang.org/x/net/ipv4"
)

type (
	PortInfo struct {
		Port uint16
		IP   IPv4
	}

	// Info represents relevant transfer information for communications.
	Info struct {
		Conn    net.Conn // optional
		IOConn  *ipv4.PacketConn
		isUDP   bool
		Local   PortInfo
		Remote  PortInfo
		Netmask IPv4
	}
)

func (self *Info) IsMulticast() bool {
	if !self.isUDP {
		return false
	}
	return self.Remote.IP.IsMulticast()
}

func (self *Info) IsUDP() bool { return self.isUDP }
func (self *Info) IsTCP() bool { return !self.isUDP }

func (self *Info) String() string {
	network := "tcp"
	if self.IsUDP() {
		network = "udp"
	}

	return fmt.Sprintf(
		"%s(%s:%X, %s:%X)",
		network,
		self.Local.IP.ToIP().String(),
		self.Local.Port,
		self.Remote.IP.ToIP().String(),
		self.Remote.Port,
	)
}

func NewDummyInfo() *Info {
	return &Info{
		isUDP: false,
		Local: PortInfo{
			Port: TCPPort,
			IP:   [4]byte{192, 160, 0, 101},
		},
		Remote: PortInfo{
			Port: TCPPort,
			IP:   [4]byte{192, 160, 0, 102},
		},
	}
}

func NewUDPInfo(
	local *net.UDPAddr,
	remote *net.UDPAddr,
	netmask IPv4,
	ioConn *ipv4.PacketConn,
) (*Info, error) {
	info := new(Info)

	info.isUDP = true
	info.IOConn = ioConn

	if err := info.Local.FromUDPAddr(local); err != nil {
		return nil, err
	}
	if err := info.Remote.FromUDPAddr(remote); err != nil {
		return nil, err
	}
	info.Netmask = netmask
	return info, nil
}

func NewTCPInfo(c net.Conn, netmask IPv4, ioConn *ipv4.PacketConn) (*Info, error) {
	if c.LocalAddr().Network() != "tcp" {
		return nil, fmt.Errorf("%s is not 'tcp'", c.LocalAddr().String())
	}
	if c.RemoteAddr().Network() != "tcp" {
		return nil, fmt.Errorf("%s is not 'tcp'", c.RemoteAddr().String())
	}

	local, err := net.ResolveTCPAddr("tcp", c.LocalAddr().String())
	if err != nil {
		return nil, err
	}
	remote, err := net.ResolveTCPAddr("tcp", c.RemoteAddr().String())
	if err != nil {
		return nil, err
	}

	info := new(Info)
	info.isUDP = false
	info.Netmask = netmask
	if err := info.Local.FromTCPAddr(local); err != nil {
		return nil, err
	}
	if err := info.Remote.FromTCPAddr(remote); err != nil {
		return nil, err
	}

	info.Conn = c
	info.IOConn = ioConn
	return info, nil
}

func (self *PortInfo) FromTCPAddr(tcpAddr *net.TCPAddr) error {
	self.Port = uint16(tcpAddr.Port)

	ipv4, err := FromIP(tcpAddr.IP)
	if err != nil {
		return err
	}
	self.IP = ipv4
	return nil
}

func (self *PortInfo) FromUDPAddr(udpAddr *net.UDPAddr) error {
	self.Port = uint16(udpAddr.Port)
	ipv4, err := FromIP(udpAddr.IP)
	if err != nil {
		return err
	}
	self.IP = ipv4
	return nil
}

func (self *Info) GetMulticastAddress() IPv4 {
	return GetMulticastAddress(self.Local.IP, self.Netmask)
}
