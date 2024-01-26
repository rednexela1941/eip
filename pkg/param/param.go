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
	onGet OnGetFunc
	onSet OnSetFunc
}

type OnGetFunc func(w bbuf.Writer) error
type OnSetFunc func(r bbuf.Reader) error

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

func _NewDefaultParam(name string, dataType DataType) *AssemblyParam {
	return &AssemblyParam{
		Name:     name,
		DataType: dataType,
	}
}

func (self *AssemblyParam) OnGet(fn OnGetFunc) *AssemblyParam {
	self.onGet = fn
	return self
}

func (self *AssemblyParam) OnSet(fn OnSetFunc) *AssemblyParam {
	self.onSet = fn
	return self
}

func NewBOOLParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, BOOL)
}

func NewSINTParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, SINT)
}

func NewINTParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, INT)
}

func NewDINTParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, DINT)
}

func NewLINTParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, LINT)
}

func NewUSINTParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, USINT)
}

func NewUINTParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, UINT)
}

func NewUDINTParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, UDINT)
}

func NewULINTParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, ULINT)
}

func NewBYTEParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, BYTE)
}

func NewWORDParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, WORD)
}

func NewDWORDParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, DWORD)
}

func NewLWORDParam(name string) *AssemblyParam {
	return _NewDefaultParam(name, LWORD)
}
