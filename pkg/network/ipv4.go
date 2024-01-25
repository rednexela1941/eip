package network

import (
	"encoding/binary"
	"fmt"
	"net"
)

// Fixed size version of IPv4 for SockaddrInfo struct
type IPv4 [4]byte

func (self IPv4) ToIP() net.IP {
	return net.IP(self[:])
}

func FromUint(v uint32) IPv4 {
	ip := IPv4{}
	for i := 0; i < len(ip); i++ {
		ip[i] = byte(0xff & (v >> (8 * (3 - i))))
	}
	return ip
}

// ToUint in big endian.
func (self IPv4) ToUint() uint32 {
	return binary.BigEndian.Uint32(self[:])
}

func (self IPv4) IsMulticast() bool {
	return 0xF0&self[0] == 0xE0
}

func (self IPv4) String() string {
	return self.ToIP().String()
}

func FromIP(ip net.IP) (IPv4, error) {
	rv := IPv4{}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return rv, fmt.Errorf("%s is not ipv4", ip.String())
	}
	for i := 0; i < 4; i++ {
		rv[i] = ipv4[i]
	}
	return rv, nil
}

func FromIPMask(mask net.IPMask) (IPv4, error) {
	nm := IPv4{}
	_, l := mask.Size()
	if l != 32 {
		return nm, fmt.Errorf("%s is not a ipv4 mask", mask.String())
	}
	if len(mask) != 4 {
		return nm, fmt.Errorf("%s is not ipv4 mask (len)", mask.String())
	}
	for i := 0; i < len(mask); i++ {
		nm[i] = mask[i]
	}
	return nm, nil
}
