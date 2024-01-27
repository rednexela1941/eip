// Package param implements the param type for ASsembly Objects
// and EDS generation.
package param

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/cip"
)

type DataType struct {
	// DataType codes for EDS.
	// See Volume 1: Table C-6.1
	Code cip.TypeCode
	Size cip.UINT
}

func (self *DataType) CodeString() string {
	return fmt.Sprintf("0x%02X", int(self.Code))
}

var (
	// BOOL = DataType{Code: BOOLCode, Size: cip.UINT(cip.BOOLSize)}
	BOOL = NewDataType(cip.BOOLCode, cip.BOOLSize)

	SINT = NewDataType(cip.SINTCode, cip.SINTSize)
	INT  = NewDataType(cip.INTCode, cip.INTSize)
	DINT = NewDataType(cip.DINTCode, cip.DINTSize)
	LINT = NewDataType(cip.LINTCode, cip.LINTSize)

	USINT = NewDataType(cip.USINTCode, cip.USINTSize)
	UINT  = NewDataType(cip.UINTCode, cip.UINTSize)
	UDINT = NewDataType(cip.UDINTCode, cip.UDINTSize)
	ULINT = NewDataType(cip.ULINTCode, cip.ULINTSize)

	BYTE  = NewDataType(cip.BYTECode, cip.BYTESize)
	WORD  = NewDataType(cip.WORDCode, cip.WORDSize)
	DWORD = NewDataType(cip.DWORDCode, cip.DWORDSize)
	LWORD = NewDataType(cip.LWORDCode, cip.LWORDSize)

	REAL  = NewDataType(cip.REALCode, cip.REALSize)
	LREAL = NewDataType(cip.LREALCode, cip.LREALSize)
)

func NewDataType(code cip.TypeCode, size int) DataType {
	return DataType{Code: code, Size: cip.UINT(size)}
}

func (self *DataType) String() string {
	return self.Code.String()
}
