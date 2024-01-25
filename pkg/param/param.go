package param

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

// AssemblyParam represents an assembly parameter.
// with data required to make a "ParamX" item in the EDS file.
// See Volume 1: 7-3.6.6.1 "Param Keyword"
type AssemblyParam struct {
	// Link Path Size USINT
	// Link Path PackedEPATH
	Index       int      // for the EDS file, shouldn't be modified directly.
	Descriptor  cip.WORD // usually 0
	DataType    DataType // type and size.
	Name        string
	UnitsString string
	HelpString  string

	MinString          string
	MaxString          string
	DefaultValueString string

	// Default Required
	onGet func(w bbuf.Writer) error
	onSet func(r bbuf.Reader) error
}

// WriteTo: for Get data (both IO connections and Get Attribute services).
func (self *AssemblyParam) WriteTo(w bbuf.Writer) error {
	if self.onGet == nil {
		return fmt.Errorf("cannot get %s", self.Name)
	}
	return self.onGet(w)
}

// ReadFrom: for Set data (both IO connection and SetAttribute services).
func (self *AssemblyParam) ReadFrom(r bbuf.Reader) error {
	if self.onSet == nil {
		return fmt.Errorf("cannot set %s", self.Name)
	}
	return self.onSet(r)
}

// SizeBits: for EDS file, equal to DataType.Size * 8
func (self *AssemblyParam) SizeBits() int {
	return int(self.DataType.Size) * 8
}

func (self *AssemblyParam) SetHelpString(s string) *AssemblyParam {
	self.HelpString = s
	return self
}

func (self *AssemblyParam) SetUnitsString(s string) *AssemblyParam {
	self.UnitsString = s
	return self
}

// SetMinString set the string to appear in the EDS file as a minimum value.
func (self *AssemblyParam) SetMinString(s string) *AssemblyParam {
	self.MinString = s
	return self
}

// SetMaxString set the string to appear in the EDS file as the maximum value.
func (self *AssemblyParam) SetMaxString(s string) *AssemblyParam {
	self.MaxString = s
	return self
}

func (self *AssemblyParam) SetDefaultValueString(s string) *AssemblyParam {
	self.DefaultValueString = s
	return self
}

func (self *AssemblyParam) GetDefaultValueString() string {
	if len(self.DefaultValueString) == 0 {
		return "0"
	}
	return self.DefaultValueString
}

func _NewDefaultParam(name string, dataType DataType, ptr any) *AssemblyParam {
	return &AssemblyParam{
		Name:     name,
		DataType: dataType,
		onGet: func(w bbuf.Writer) error {
			w.Wl(ptr)
			return w.Error()
		},
		onSet: func(r bbuf.Reader) error {
			r.Rl(ptr)
			return r.Error()
		},
	}
}

func NewBOOLParam(name string, ptr *cip.BOOL) *AssemblyParam {
	return _NewDefaultParam(name, BOOL, ptr)
}

func NewSINTParam(name string, ptr *cip.SINT) *AssemblyParam {
	return _NewDefaultParam(name, SINT, ptr)
}

func NewINTParam(name string, ptr *cip.INT) *AssemblyParam {
	return _NewDefaultParam(name, INT, ptr)
}

func NewDINTParam(name string, ptr *cip.DINT) *AssemblyParam {
	return _NewDefaultParam(name, DINT, ptr)
}

func NewLINTParam(name string, ptr *cip.LINT) *AssemblyParam {
	return _NewDefaultParam(name, LINT, ptr)
}

func NewUSINTParam(name string, ptr *cip.USINT) *AssemblyParam {
	return _NewDefaultParam(name, USINT, ptr)
}

func NewUINTParam(name string, ptr *cip.UINT) *AssemblyParam {
	return _NewDefaultParam(name, UINT, ptr)
}

func NewUDINTParam(name string, ptr *cip.UDINT) *AssemblyParam {
	return _NewDefaultParam(name, UDINT, ptr)
}

func NewULINTParam(name string, ptr *cip.ULINT) *AssemblyParam {
	return _NewDefaultParam(name, ULINT, ptr)
}

func NewBYTEParam(name string, ptr *cip.BYTE) *AssemblyParam {
	return _NewDefaultParam(name, BYTE, ptr)
}

func NewWORDParam(name string, ptr *cip.WORD) *AssemblyParam {
	return _NewDefaultParam(name, WORD, ptr)
}

func NewDWORDParam(name string, ptr *cip.DWORD) *AssemblyParam {
	return _NewDefaultParam(name, DWORD, ptr)
}

func NewLWORDParam(name string, ptr *cip.LWORD) *AssemblyParam {
	return _NewDefaultParam(name, LWORD, ptr)
}
