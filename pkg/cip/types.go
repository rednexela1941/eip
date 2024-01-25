// Package cip provides various types and constant defintions related to Common Industrial Protocol.
package cip

// See Vol 1: Appendix C-2.1.1 "Elementary Data Types"
type (
	OCTET        = byte
	BOOL         = bool
	SINT         = int8
	INT          = int16
	DINT         = int32
	LINT         = int64
	USINT        = uint8
	UINT         = uint16
	UDINT        = uint32
	ULINT        = uint64
	REAL         = float32
	LREAL        = float64
	BYTE         = uint8
	WORD         = uint16
	DWORD        = uint32
	LWORD        = uint64
	SHORT_STRING = []byte // []{length, char0, char1, ..., charN}
)

const (
	OCTETSize int = 1 // bytes
	BOOLSize      = 1
	SINTSize      = 1
	INTSize       = 2
	DINTSize      = 4
	LINTSize      = 8
	USINTSize     = 1
	UINTSize      = 2
	UDINTSize     = 4
	ULINTSize     = 8
	REALSize      = 4
	LREALSize     = 8
	BYTESize      = 1
	WORDSize      = 2
	DWORDSize     = 4
	LWORDSize     = 8
)
