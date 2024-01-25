package encap

import "github.com/rednexela1941/eip/pkg/cip"

// See: Volume 2-3.1

//go:generate stringer -type=Command
type Command cip.UINT

// See Volume 2: 2-3.2 (Command Field)
const (
	NOP               Command = 0x0000 // TCP
	ListServices      Command = 0x0004
	ListIdentity      Command = 0x0063 // see Vol2: 2-4.2.3 for reply rules (random delay, etc.)
	ListInterfaces    Command = 0x0064
	RegisterSession   Command = 0x0065 // TCP
	UnRegisterSession Command = 0x0066 // TCP
	SendRRData        Command = 0x006F // "RR" = RequestReply.
	SendUnitData      Command = 0x0070 // ony allowed TCP.
	StartDTLS         Command = 0x00C8 // no TLS right now.
	// ReservedForLegacyUsage = 0x0005
	// ReservedForFutureExpansion = 0x0006 - 0x0062
	// ReservedForLegacyUsage = 0x0001-0x0003
	// ReservedForFutureExpansion2 = 0x00C9 - 0xFFFF
	// ReservedForLegacyUsage = 0x0071 - 0x00C7
)
