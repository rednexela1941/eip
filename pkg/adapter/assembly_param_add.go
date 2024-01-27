package adapter

import (
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/param"
)

func (self *AssemblyInstance) AddBOOLParam(name string) *ElementaryParam[cip.BOOL] {
	ep := _NewElementaryParam[cip.BOOL](name, param.BOOL)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddSINTParam(name string) *ElementaryParam[cip.SINT] {
	ep := _NewElementaryParam[cip.SINT](name, param.SINT)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddINTParam(name string) *ElementaryParam[cip.INT] {
	ep := _NewElementaryParam[cip.INT](name, param.INT)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddDINTParam(name string) *ElementaryParam[cip.DINT] {
	ep := _NewElementaryParam[cip.DINT](name, param.DINT)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddLINTParam(name string) *ElementaryParam[cip.LINT] {
	ep := _NewElementaryParam[cip.LINT](name, param.LINT)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddUSINTParam(name string) *ElementaryParam[cip.USINT] {
	ep := _NewElementaryParam[cip.USINT](name, param.USINT)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddUINTParam(name string) *ElementaryParam[cip.UINT] {
	ep := _NewElementaryParam[cip.UINT](name, param.UINT)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddUDINTParam(name string) *ElementaryParam[cip.UDINT] {
	ep := _NewElementaryParam[cip.UDINT](name, param.UDINT)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddULINTParam(name string) *ElementaryParam[cip.ULINT] {
	ep := _NewElementaryParam[cip.ULINT](name, param.ULINT)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddBYTEParam(name string) *ElementaryParam[cip.BYTE] {
	ep := _NewElementaryParam[cip.BYTE](name, param.BYTE)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddWORDParam(name string) *ElementaryParam[cip.WORD] {
	ep := _NewElementaryParam[cip.WORD](name, param.WORD)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddDWORDParam(name string) *ElementaryParam[cip.DWORD] {
	ep := _NewElementaryParam[cip.DWORD](name, param.DWORD)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddLWORDParam(name string) *ElementaryParam[cip.LWORD] {
	ep := _NewElementaryParam[cip.LWORD](name, param.LWORD)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddREALParam(name string) *ElementaryParam[cip.REAL] {
	ep := _NewElementaryParam[cip.REAL](name, param.REAL)
	self.AddParam(ep.AssemblyParam)
	return ep
}

func (self *AssemblyInstance) AddLREALParam(name string) *ElementaryParam[cip.LREAL] {
	ep := _NewElementaryParam[cip.LREAL](name, param.LREAL)
	self.AddParam(ep.AssemblyParam)
	return ep
}
