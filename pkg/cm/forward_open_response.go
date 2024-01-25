package cm

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

type ForwardOpenResponseHeader struct {
	// TODO: note rules for choosing connection ids.
	OtoTNetworkConnectionID   cip.UDINT
	TtoONetworkConnectionID   cip.UDINT
	Triad                               // Same as request packet.
	OtoTAPI                   cip.UDINT // Actual Packet Interval (microseconds)
	TtoOAPI                   cip.UDINT // Actual Packet Interval (microseconds)
	ApplicationReplySizeWords cip.USINT // number of words in application reply.
	Reserved                  cip.USINT
}

// Volume 1: Table 3-5.22
type ForwardOpenResponse struct {
	ForwardOpenResponseHeader // seperate out (for easier encoding).
	ApplicationReply          []cip.WORD
}

type RemainingPathSize struct {
	Words    cip.USINT
	Reserved cip.OCTET
}

// Volume 1: Table 3-5.23
type UnsuccessfulForwardOpenResponse struct {
	Triad
	// RemainingPathSize is optional.
	RemainingPathSize
}

func WriteUnsuccessfulForwardOpenResponse(w bbuf.Writer, freq *SharedForwardOpenRequest) {
	w.Wl(&freq.Triad)
}
