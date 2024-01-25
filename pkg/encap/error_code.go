package encap

import "github.com/rednexela1941/eip/pkg/cip"

//go:generate stringer -type=ErrorCode
type ErrorCode cip.UDINT // Vol1 2-3.5 (for status field in encap header)
// See Vol 2: Table 2-3.3
const (
	StatusSuccess                                    ErrorCode = 0x0000
	StatusInvalidCommand                             ErrorCode = 0x0001
	StatusInsufficientMemory                         ErrorCode = 0x0002
	StatusInvalidSessionHandle                       ErrorCode = 0x0064
	StatusInvalidLength                              ErrorCode = 0x0065
	StatusUnsupportedProtocolVersion                 ErrorCode = 0x0069
	StatusEncapsulatedCIPServiceNotAllowedOnThisPort ErrorCode = 0x006A
)
