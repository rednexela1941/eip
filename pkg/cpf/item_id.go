package cpf

import "github.com/rednexela1941/eip/pkg/cip"

//go:generate stringer -type=ItemID
type ItemID cip.UINT

const (
	// CPF Item ID's.
	// see Vol2 : Table 2-5.3
	// and Vol2 : 2-5.2 for usage/definitions.
	AddressNull               ItemID = 0x0
	CIPIdentity               ItemID = 0xC
	SecurityInfo              ItemID = 0x86
	EtherNetIPCapability      ItemID = 0x87
	EtherNetIPUsage           ItemID = 0x88
	AddressConnected          ItemID = 0xA1
	ConnectedTransportPacket  ItemID = 0xB1
	UnconnectedMessage        ItemID = 0xB2
	ListServicesResponse      ItemID = 0x100 // Vol2 : 2-4.6.3
	SockaddrInfoOtoT          ItemID = 0x8000
	SockaddrInfoTtoO          ItemID = 0x8001
	SequencedAddress          ItemID = 0x8002
	UnconnectedMessageOverUDP ItemID = 0x8003 // Vol 8: Chapter 3
)
