package network

import "net"

// Volume 2: 3-5.3 (Multicast Address Allocation)
// TODO: Multicast Algo. (Vol2 : 3-5.3)
var MulticastBaseAddress = net.IPv4(239, 192, 1, 0)

//	func GetMulticastAddress() IPv4 {
//		// TODO: make this use appropriate address somewhere else.
//		return _GetMulticastAddress(
//			IPv4{192, 168, 0, 101},
//			IPv4{255, 255, 255, 0},
//		)
//	}
//
// See Volume 2: 3-5.3 Multicast Address Allocation for EtherNet/IP.
func GetMulticastAddress(ipAddr IPv4, netmask IPv4) IPv4 {
	mcast := IPv4{239, 192, 1, 0}.ToUint()
	ip := ipAddr.ToUint()
	nm := netmask.ToUint()

	hostID := ip & ^nm
	mcastIndex := hostID - 1
	mcastIndex &= 0x3FF // 10 bits of host ID.
	startAddr := mcast + (mcastIndex * 32)

	return FromUint(startAddr)
}
