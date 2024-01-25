// Package encap implements types and functionality for Encapsulation Messages (Volume 2: 2-3)
package encap

import (
	"github.com/rednexela1941/eip/pkg/cip"
)

// See Vol 2: 2-3
const HeaderLength = 24

type (
	SessionHandle = cip.UDINT
	SenderContext = [8]cip.OCTET

	// See Volume 2: Table 2-3.1
	Packet struct {
		Header
		CommandSpecificData []cip.OCTET
	}

	Header struct {
		Command       Command
		Length        cip.UINT
		SessionHandle SessionHandle /* UDINT */
		Status        ErrorCode
		SenderContext SenderContext // only relevant to sender, length 8
		Options       cip.UDINT
	}
)
