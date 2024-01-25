package cpf

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/network"
)

// See Vol 2: 2-5.3.3
// Sockaddr Info Item (Vol 2: 2-5.3.3, Table 2-5.9). Note that all fields encoded as big-endian.

type (
	SockaddrInfo struct {
		// Table 2-5.9 in volume 2

		// shall be AF_INET = 2. This field shall be sent in big endian order.
		SinFamily cip.INT
		// For point-point connections, sin_port shall be set to the UDP port to which packets for this CIP connection will be sent. For point-point connections, it is recommended that the registered UDP port (0x8AE) be used. When used with a multicast connection, the sin_port field shall be set to the registered UDP port number (0x08AE) and treated by the receiver as “don’tcare”. This field shall be sent in big endian order.
		SinPort cip.UINT // UDP port number
		// For multicast connections, sin_addr shall be set to the IP multicast address to which packets for this CIP connection will be sent.  When used with a point-point connection, the sin_addr field shall be treated by the receiver as “don’t care”.  It is recommended that the sender set sin_addr to 0 for point-point connections. This field shall be sent in big endian order.
		SinAddr network.IPv4
		// Length of 8. Recommended value of zero; not enforced
		SinZero [8]cip.USINT // all zero.
	}
)

// NB: Also see Vol 2: Table 2-5.10 (forward open behavior, must send sockaddr info.)
// NB: Also see Vol 2: Table 3-3.2 (Sockaddr Info Usage for Forward Open service)

func NewSockaddrInfo(ip network.IPv4, port uint16) SockaddrInfo {
	return SockaddrInfo{
		SinFamily: 2,
		SinAddr:   ip,
		SinPort:   port,
	}
}

func (self *SockaddrInfo) String() string {
	return fmt.Sprintf("%s:%d", self.SinAddr.String(), self.SinPort)
}
