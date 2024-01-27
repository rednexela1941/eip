package adapter

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/param"
)

// ElementaryParam represents an assembly parameter
// of one of the cip.ELEMENTARY types.
type ElementaryParam[T cip.ELEMENTARY] struct {
	*param.AssemblyParam

	onGet func() T
	onSet func(T)
}

func _NewElementaryParam[T cip.ELEMENTARY](
	name string,
	dataType param.DataType,
) *ElementaryParam[T] {
	ep := &ElementaryParam[T]{
		AssemblyParam: param.NewDefaultParam(name, dataType),
	}
	ep.AssemblyParam.OnGet(func(w bbuf.Writer) error {
		return ep.get(w)
	})
	ep.AssemblyParam.OnSet(func(r bbuf.Reader) error {
		return ep.set(r)
	})
	return ep
}

func (self *ElementaryParam[T]) OnGet(fn func() T) *ElementaryParam[T] {
	self.onGet = fn
	return self
}

func (self *ElementaryParam[T]) OnSet(fn func(T)) *ElementaryParam[T] {
	self.onSet = fn
	return self
}

func (self *ElementaryParam[T]) SetHelpString(s string) *ElementaryParam[T] {
	self.HelpString = s
	return self
}

func (self *ElementaryParam[T]) SetUnitsString(s string) *ElementaryParam[T] {
	self.UnitsString = s
	return self
}

// SetMinString set the string to appear in the EDS file as a minimum value.
func (self *ElementaryParam[T]) SetMinString(s string) *ElementaryParam[T] {
	self.MinString = s
	return self
}

// SetMaxString set the string to appear in the EDS file as the maximum value.
func (self *ElementaryParam[T]) SetMaxString(s string) *ElementaryParam[T] {
	self.MaxString = s
	return self
}

func (self *ElementaryParam[T]) SetDefaultValueString(s string) *ElementaryParam[T] {
	self.DefaultValueString = s
	return self
}

func (self *ElementaryParam[T]) get(w bbuf.Writer) error {
	if self.onGet == nil {
		return fmt.Errorf("cannot get %s", self.Name)
	}
	v := self.onGet()
	w.Wl(&v)
	return w.Error()
}

func (self *ElementaryParam[T]) set(r bbuf.Reader) error {
	if self.onSet == nil {
		return fmt.Errorf("cannot set %s", self.Name)
	}
	var v T
	r.Rl(&v)
	if r.Error() != nil {
		return r.Error()
	}
	self.onSet(v)
	return nil
}
