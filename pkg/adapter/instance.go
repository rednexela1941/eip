package adapter

import (
	"fmt"
	"log"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/epath"
)

type (
	AttributeMap map[cip.UINT]*Attribute

	Instance struct {
		ServiceMap
		ClassCode  cip.ClassCode
		InstanceID cip.UINT
		Attributes AttributeMap
	}
)

func NewInstance(classCode cip.ClassCode, instanceID cip.UINT) *Instance {
	return &Instance{
		ServiceMap: make(ServiceMap),
		ClassCode:  classCode,
		InstanceID: instanceID,
		Attributes: make(AttributeMap),
	}
}

func (self *Instance) GetPath() *epath.ApplicationPath {
	ap := new(epath.ApplicationPath)
	ap.ClassID = self.ClassCode
	ap.InstanceID = self.InstanceID
	return ap
}

func (self *Instance) CanCall(serviceCode cip.ServiceCode) bool {
	_, ok := self.ServiceMap[serviceCode]
	if ok {
		return true
	}

	for _, attr := range self.Attributes {
		if _, ok := attr.ServiceMap[serviceCode]; ok {
			return true
		}
	}
	return false
}

func (self *Instance) ClassCodeHex() string {
	return fmt.Sprintf("0x%X", int(self.ClassCode))
}

func (self *Instance) CallService(req *Request, res Response) {
	if !self.CanCall(req.Service) {
		res.SetGeneralStatus(cip.StatusServiceNotSupported)
		return
	}

	path := req.RequestPath.ApplicationPaths[0]

	if path.AttributeID == 0 {
		// call directly on instance.
		serviceFn, ok := self.ServiceMap[req.Service]
		if !ok {
			// may have to rearrange when this happens.
			res.SetGeneralStatus(cip.StatusAttributeNotSupported)
			// res.SetGeneralStatus(cip.StatusServiceNotSupported)
			return
		}
		serviceFn(req, res)
		return
	}

	attribute, ok := self.Attributes[path.AttributeID]
	if !ok {
		res.SetGeneralStatus(cip.StatusAttributeNotSupported)
		return
	}
	attribute.CallService(req, res)
}

func (self *Instance) AddAttribute(attrID cip.UINT, name string, size int) *Attribute {
	a := NewAttribute(attrID, name, size)
	self.Attributes[attrID] = a
	return a
}

func (self *Instance) HighestAttributeID() cip.UINT {
	maxID := cip.UINT(0)
	for id := range self.Attributes {
		if id > maxID {
			maxID = id
		}
	}
	return maxID
}

func (self *Instance) addDefaultGetAttributesAll() {
	// default setup for get attributes all.
	// which uses attributes to form response.
	fn := func(req *Request, res Response) {
		didCall := false
		for i := cip.UINT(0); i <= self.HighestAttributeID(); i++ {
			attribute, ok := self.Attributes[i]
			if !ok {
				continue
			}
			serviceFn, ok := attribute.ServiceMap[cip.GetAttributesAll]
			if !ok {
				continue
			}
			serviceFn(req, res)
			didCall = true
		}
		if !didCall {
			// if no attributes support GetAttributesAll
			res.SetGeneralStatus(cip.StatusServiceNotSupported)
		}
	}

	self.OnService(cip.GetAttributesAll, fn)
}

func (self *Instance) addDefaultSetAttributeList() {
	// add default handler for setattributes list.
	// See Volume 1: A-4.4.2
	// And Volume 1: Table A-4.5
	// And Volume 1: Table A-4.6

	fn := func(req *Request, res Response) {
		r := bbuf.New(req.RequestData)
		var numAttrs cip.UINT
		r.Rl(&numAttrs)
		res.Wl(numAttrs)

		for i := cip.UINT(0); i < numAttrs; i++ {
			var attrID cip.UINT
			r.Rl(&attrID)
			res.Wl(attrID)

			attribute, ok := self.Attributes[attrID]
			if !ok {
				log.Printf("attribute %d not found\n", attrID)
				res.SetGeneralStatus(cip.StatusAttributeListError)
				res.Wl(cip.StatusAttributeNotSupported)
				res.Wl(cip.OCTET(0)) // reserved
				return
			}

			data := make([]cip.BYTE, attribute.Size)
			_, err := r.Read(data)
			if err != nil {
				res.SetGeneralStatus(cip.StatusNotEnoughData)
				log.Println(err)
				return
			}

			serviceFn, ok := attribute.ServiceMap[cip.SetAttributeList]
			if !ok {
				log.Printf("%s not supported on %s\n",
					req.Service.String(),
					attribute.Name,
				)
				res.SetGeneralStatus(cip.StatusAttributeListError)
				res.Wl(cip.StatusAttributeNotSetable)
				res.Wl(cip.OCTET(0)) // reserved
				continue
			}

			res.Wl(cip.StatusSuccess)
			res.Wl(cip.OCTET(0)) // reseved.
			serviceFn(req, res)
		}
	}
	self.OnService(cip.SetAttributeList, fn)
}

func (self *Instance) addDefaultGetAttributesList() {
	// add default service handler for getattributes list.
	// using self.Attributes.

	// See Volume 1: A-4.3.3
	// and Volume 1: Table A-4.4
	fn := func(req *Request, res Response) {
		r := bbuf.New(req.RequestData)
		var numAttrs cip.UINT
		r.Rl(&numAttrs)
		res.Wl(numAttrs)

		for i := cip.UINT(0); i < numAttrs; i++ {
			var attrID cip.UINT
			r.Rl(&attrID)

			res.Wl(attrID)

			attribute, ok := self.Attributes[attrID]
			if !ok {
				log.Printf("attribute %d not found\n", attrID)
				res.SetGeneralStatus(cip.StatusAttributeListError)
				res.Wl(cip.StatusAttributeNotSupported)
				res.Wl(cip.OCTET(0)) // reserved
				continue
			}
			serviceFn, ok := attribute.ServiceMap[cip.GetAttributeList]
			if !ok {
				log.Printf("%s not supported on %s\n",
					req.Service.String(),
					attribute.Name,
				)
				res.SetGeneralStatus(cip.StatusAttributeListError)
				res.Wl(cip.StatusServiceNotSupported)
				res.Wl(cip.OCTET(0)) // reserved
				continue
			}

			res.Wl(cip.StatusSuccess)
			res.Wl(cip.OCTET(0)) // reseved.
			serviceFn(req, res)
		}

		if r.Error() != nil {
			// bad packet.
			res.SetGeneralStatus(cip.StatusNotEnoughData)
		}
	}

	self.OnService(cip.GetAttributeList, fn)
}
