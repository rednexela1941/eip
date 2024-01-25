package cm

import "github.com/rednexela1941/eip/pkg/cip"

// Volume 1: Table 3-5.25
type ForwardCloseResponse struct {
	Triad
	ApplicationReplySizeWords cip.USINT
	Reserved                  cip.USINT
	ApplicationReply          []cip.WORD
}

// Volume 1: Table 3-5.26
type UnsuccessfulForwardCloseResponse struct {
	Triad
	RemainingPathSize // optional, see Volume 1 Table 3-5.26
}
