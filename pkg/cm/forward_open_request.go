package cm

import (
	"time"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/epath"
)

const (
	ForwardOpenHeaderLength      int = 36
	LargeForwardOpenHeaderLength     = 40
)

type (
	// Volume 1: 3-5.6.1.7
	Triad struct {
		ConnectionSerialNumber cip.UINT
		OriginatorVendorID     cip.UINT
		OriginatorSerialNumber cip.UDINT
	}

	// ForwardOpenRequestHeader is shared between LargeForwardOpenRequests and ForwardOpenRequests.
	ForwardOpenRequestHeader struct {
		PriorityTimeTick PriorityTimeTick // see 3-5.6.1.2.1
		TimeoutTicks     cip.USINT        // see 3-5.6.1.2.1

		// See Volume 2: Table 3-3.2 Network Connection ID Selection
		OtoTNetworkConnectionID cip.UDINT
		TtoONetworkConnectionID cip.UDINT
		// Volume 1: 3-5.6.1.7
		Triad
		// Volume 1: 3-5.6.1.4 "Connection Timeout Multiplier"
		// and Volume 1: Table 3-5.15
		ConnectionTimeoutMultiplier cip.USINT
		Reserved                    [3]cip.OCTET
	}

	// See Volume 1: Table 3-5.11
	_ConnectionParameters struct {
		RPI        cip.UDINT // microseconds
		Parameters cip.WORD
	}

	// See Volume 1: Table 3-5.12
	_LargeConnectionParameters struct {
		RPI        cip.UDINT // microseconds
		Parameters cip.DWORD
	}

	// See Volume 1: Table 3-5.20
	_ForwardOpenRequest struct {
		ForwardOpenRequestHeader
		// Microseconds
		OtoTParameters _ConnectionParameters
		TtoOParameters _ConnectionParameters
		ForwardOpenRequestFooter
	}

	// See Volume 1: Table 3-5.20
	_LargeForwardOpenRequest struct {
		ForwardOpenRequestHeader
		OtoTParameters _LargeConnectionParameters
		TtoOParameters _LargeConnectionParameters
		ForwardOpenRequestFooter
	}

	// ForwardOpenRequestFooter is shared between LargeForwardOpen and ForwardOpen.
	ForwardOpenRequestFooter struct {
		TransportClassAndTrigger TransportClassAndTrigger
		ConnectionPathSizeWords  cip.USINT
	}

	// See Volume 1: Table 3-5.20
	// Represents both ForwardOpen and LargeForwardOpen requests.
	// struct is "verbose" -- mean to be quickly readable when deubgging at
	// the cost of some extra size.
	SharedForwardOpenRequest struct {
		ForwardOpenRequestHeader
		OtoTParameters Parameters
		TtoOParameters Parameters
		ForwardOpenRequestFooter
		// See Volume 1: Table 3-5.16
		// And Volume 1: 3-5.6.1.10
		ConnectionPath *epath.PaddedEPATH
	}
)

// IsNull (is Null Forward Open Request)
func (self *SharedForwardOpenRequest) IsNull() bool {
	return self.TtoOParameters.IsNull() && self.OtoTParameters.IsNull()
}

// See Volume 1: 3-5.6.1.2.1 "Unconnected Request Timing"
func (self *ForwardOpenRequestHeader) GetTimeout() time.Duration {
	return self.PriorityTimeTick.TickTime() * time.Duration(self.TimeoutTicks)
}

// See Volume 1: Table 3-5.15 "Connection Timeout Mulitplier Values"
// this will be applied to the RPI to determine connection timeout values.
// returns (multiplier int, ok bool)
func (self *ForwardOpenRequestHeader) GetTimeoutMultiplier() (int, bool) {
	v := self.ConnectionTimeoutMultiplier
	if v >= 8 {
		return 0, false
	}
	return 1 << (int(v) + 2), true
}

func _NewForwardOpenRequest(r bbuf.Reader, params intoAble) *SharedForwardOpenRequest {
	f := new(SharedForwardOpenRequest)

	r.Rl(&f.ForwardOpenRequestHeader)

	r.Rl(params) // OtoT Parameters
	f.OtoTParameters = params.Into()

	r.Rl(params) // TtoO Parameters
	f.TtoOParameters = params.Into()

	r.Rl(&f.ForwardOpenRequestFooter)

	if f.ConnectionPathSizeWords > 0 {
		cpData := make([]byte, int(f.ConnectionPathSizeWords)*2)
		r.Rl(&cpData)
		f.ConnectionPath = epath.NewPadded(cpData)
	}

	return f
}

// NewForwardOpen request parses reader into a SharedFowardOpenRequest structure.
// See Volume 1: Table 3-5.20
// Check for reader errors after calling.
func NewForwardOpenRequest(r bbuf.Reader) *SharedForwardOpenRequest {
	params := new(_ConnectionParameters)
	return _NewForwardOpenRequest(r, params)
}

// NewLargeForwardOpen request parses reader into a SharedFowardOpenRequest structure.
// See Volume 1: Table 3-5.20
// Check for reader errors after calling.
func NewLargeForwardOpenRequest(r bbuf.Reader) *SharedForwardOpenRequest {
	params := new(_LargeConnectionParameters)
	return _NewForwardOpenRequest(r, params)
}
