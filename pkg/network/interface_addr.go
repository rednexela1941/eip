package network

import "net"

type InterfaceAddr struct {
	IP      IPv4
	Netmask IPv4
}

func FromIPNet(ipnet *net.IPNet) (*InterfaceAddr, error) {
	ip, err := FromIP(ipnet.IP)
	if err != nil {
		return nil, err
	}
	mask, err := FromIPMask(ipnet.Mask)
	if err != nil {
		return nil, err
	}
	return &InterfaceAddr{
		IP:      ip,
		Netmask: mask,
	}, nil
}

func (self *InterfaceAddr) GetMulticastAddress() IPv4 {
	return GetMulticastAddress(self.IP, self.Netmask)
}
