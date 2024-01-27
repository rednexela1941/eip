// Package cip provides various types and constant defintions related to Common Industrial Protocol.
package cip

import (
	"golang.org/x/exp/constraints"
)

// See Vol 1: Appendix C-2.1.1 "Elementary Data Types"
type (
	OCTET = byte
	BOOL  = bool
	SINT  = int8
	INT   = int16
	DINT  = int32
	LINT  = int64
	USINT = uint8
	UINT  = uint16
	UDINT = uint32
	ULINT = uint64
	REAL  = float32
	LREAL = float64
	BYTE  = uint8
	WORD  = uint16
	DWORD = uint32
	LWORD = uint64

	// short string
	SHORT_STRING = []byte // []{length, char0, char1, ..., charN}

	ELEMENTARY interface {
		// any elementary type.
		constraints.Integer | constraints.Float | ~bool
	}
)

// DataType codes for EDS.
// See Volume 1: Table C-6.1
type TypeCode USINT

//go:generate stringer -type=TypeCode

const (
	UTIMECode TypeCode = 0xC0
	BOOLCode  TypeCode = 0xC1

	SINTCode TypeCode = 0xC2
	INTCode  TypeCode = 0xC3
	DINTCode TypeCode = 0xC4
	LINTCode TypeCode = 0xC5

	USINTCode TypeCode = 0xC6
	UINTCode  TypeCode = 0xC7
	UDINTCode TypeCode = 0xC8
	ULINTCode TypeCode = 0xC9

	REALCode  TypeCode = 0xCA
	LREALCode TypeCode = 0xCB

	BYTECode  TypeCode = 0xD1
	WORDCode  TypeCode = 0xD2
	DWORDCode TypeCode = 0xD3
	LWORDCode TypeCode = 0xD4
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
