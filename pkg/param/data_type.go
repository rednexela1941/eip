// Package param implements the param type for ASsembly Objects
// and EDS generation.
package param

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/cip"
)

//go:generate stringer -type=DataTypeCode

// DataType codes for EDS.
// See Volume 1: Table C-6.1
type DataTypeCode cip.USINT

type DataType struct {
	Code DataTypeCode
	Size cip.UINT
}

func (self *DataType) CodeString() string {
	return fmt.Sprintf("0x%02X", int(self.Code))
}

const (
	UTIMECode DataTypeCode = 0xC0
	BOOLCode  DataTypeCode = 0xC1

	SINTCode DataTypeCode = 0xC2
	INTCode  DataTypeCode = 0xC3
	DINTCode DataTypeCode = 0xC4
	LINTCode DataTypeCode = 0xC5

	USINTCode DataTypeCode = 0xC6
	UINTCode  DataTypeCode = 0xC7
	UDINTCode DataTypeCode = 0xC8
	ULINTCode DataTypeCode = 0xC9

	REALCode  DataTypeCode = 0xCA
	LREALCode DataTypeCode = 0xCB

	BYTECode  DataTypeCode = 0xD1
	WORDCode  DataTypeCode = 0xD2
	DWORDCode DataTypeCode = 0xD3
	LWORDCode DataTypeCode = 0xD4
)

var (
	// BOOL = DataType{Code: BOOLCode, Size: cip.UINT(cip.BOOLSize)}
	BOOL = NewDataType(BOOLCode, cip.BOOLSize)

	SINT = NewDataType(SINTCode, cip.SINTSize)
	INT  = NewDataType(INTCode, cip.INTSize)
	DINT = NewDataType(DINTCode, cip.DINTSize)
	LINT = NewDataType(LINTCode, cip.LINTSize)

	USINT = NewDataType(USINTCode, cip.USINTSize)
	UINT  = NewDataType(UINTCode, cip.UINTSize)
	UDINT = NewDataType(UDINTCode, cip.UDINTSize)
	ULINT = NewDataType(ULINTCode, cip.ULINTSize)

	BYTE  = NewDataType(BYTECode, cip.BYTESize)
	WORD  = NewDataType(WORDCode, cip.WORDSize)
	DWORD = NewDataType(DWORDCode, cip.DWORDSize)
	LWORD = NewDataType(LWORDCode, cip.LWORDSize)
)

func NewDataType(code DataTypeCode, size int) DataType {
	return DataType{Code: code, Size: cip.UINT(size)}
}

func (self *DataType) String() string {
	return self.Code.String()
}
