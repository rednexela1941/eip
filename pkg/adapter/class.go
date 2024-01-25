package adapter

import (
	"github.com/rednexela1941/eip/pkg/cip"
)

type (
	InstanceMap map[cip.UINT]*Instance
	ClassMap    map[cip.ClassCode]*Class

	Class struct {
		Name     string
		Revision cip.UINT
		*Instance
		Instances InstanceMap
	}
)

func NewClass(name string, classCode cip.ClassCode, revision cip.UINT) *Class {
	instance := NewInstance(classCode, 0 /* class instance is ID zero */)
	class := new(Class)
	class.Instance = instance
	class.Instances = make(InstanceMap)
	class.Name = name
	class.Revision = revision
	return class
}

func (self *Class) CallService(req *Request, res Response) {
	path := req.RequestPath.ApplicationPaths[0]

	if path.InstanceID == 0 {
		// targets the class itself.
		self.Instance.CallService(req, res)
		return
	}

	instance, ok := self.Instances[path.InstanceID]
	if !ok {
		res.SetGeneralStatus(cip.StatusObjectInstanceDoesNotExist)
		return
	}

	instance.CallService(req, res)
	return
}

func (self *Class) AddInstance(instanceID cip.UINT) *Instance {
	instance := NewInstance(self.ClassCode, instanceID)
	self.Instances[instanceID] = instance
	return instance
}

func (self *Class) NumberOfInstances() cip.UINT {
	return cip.UINT(len(self.Instances))
}

func (self *Class) HighestInstanceID() cip.UINT {
	maxID := cip.UINT(0)
	for id := range self.Instances {
		if id > maxID {
			maxID = id
		}
	}
	return maxID
}

func (self *Class) HighestInstanceAttributeID() cip.UINT {
	maxID := cip.UINT(0)
	for _, instance := range self.Instances {
		id := instance.HighestAttributeID()
		if id > maxID {
			maxID = id
		}
	}
	return maxID
}
