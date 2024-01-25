package adapter

import (
	"github.com/rednexela1941/eip/pkg/cip"
)

func (self *_Adapter) AddClass(name string, classCode cip.ClassCode, revision cip.UINT) *Class {
	c := NewClass(name, classCode, revision)
	self.Classes[classCode] = c
	return c
}
