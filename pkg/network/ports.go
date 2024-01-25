// Package network provides EtherNet/IP related networking constants (port numbers) and interface definitions.
package network

const (
	// Vol 2: Table 2-2.1
	TCPPort    uint16 = 0xAF12
	TCPTLSPort uint16 = 0x08AD
	// Vol 2: Table 2-2.2
	UDPPort     uint16 = 0xAF12
	UDPIOPort   uint16 = 0x08AE // see Vol1: Chapter 3
	UDPDTLSPort uint16 = 0x08AD
)
