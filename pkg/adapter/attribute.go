package adapter

import (
	"fmt"
	"log"

	"github.com/rednexela1941/eip/pkg/cip"
)

const (
	GetSingle AttributeFlags = 1 << iota
	GetAll
	GetList
	SetSingle
	SetAll
	SetList

	GetFull = GetSingle | GetAll | GetList
)

// Volume 1 : Table 4-4.2
const (
	AttributeClassRevision                     cip.UINT = 1
	AttributeMaxInstance                       cip.UINT = 2
	AttributeNumberOfInstances                 cip.UINT = 3
	AttributeOptionalAttributeList             cip.UINT = 4
	AttributeOptionalServiceList               cip.UINT = 5
	AttributeMaximumIDNumberClassAttributes    cip.UINT = 6
	AttributeMaximumIDNumberInstanceAttributes cip.UINT = 7
)

type (
	AttributeFlags = int

	// Attribute represents simple get/set attributes
	// that are commonly used in class/instances defintions.
	// TODO: How should we organize values?
	Attribute struct {
		Name string // for debugging purposes
		ServiceMap
		AttributeID cip.UINT
		Size        int // in bytes, TODO: use in set functions to check for enough data.
		// Flags       AttributeFlags
	}

	// ServiceHandler represents a generic CIP service call handler.
	ServiceHandler func(*Request, Response)

	ServiceMap map[cip.ServiceCode]ServiceHandler

	// GetEncoderFunc represents a function for encoding GetAttribute(Single|All|List) requests on an attribute
	GetEncoderFunc func(res Response)
)

func NewAttribute(attrID cip.UINT, name string, size int) *Attribute {
	return &Attribute{
		Name:        name,
		AttributeID: attrID,
		ServiceMap:  make(ServiceMap),
		Size:        size,
		// Flags:       flags,
	}
}

func (self *Attribute) callGetSingle(req *Request, res Response) error {
	serviceFn, ok := self.ServiceMap[cip.GetAttributeSingle]
	if !ok {
		return fmt.Errorf("service not supported")
	}
	serviceFn(req, res)
	return nil
}

func (self *Attribute) CallService(req *Request, res Response) {
	serviceFn, ok := self.ServiceMap[req.Service]
	if !ok {
		if req.Service.IsGet() {
			res.SetGeneralStatus(cip.StatusAttributeNotGettable)
			return
		}
		if req.Service.IsSet() {
			res.SetGeneralStatus(cip.StatusAttributeNotSetable)
			return
		}
		res.SetGeneralStatus(cip.StatusServiceNotSupported)
		return
	}
	serviceFn(req, res)
}

func (self ServiceMap) OnService(serviceCode cip.ServiceCode, handler ServiceHandler) {
	self[serviceCode] = handler
}

func (self *Attribute) OnSet(flags AttributeFlags, fn ServiceHandler) *Attribute {
	if flags&SetSingle != 0 {
		self.OnService(cip.SetAttributeSingle, fn)
	}
	if flags&SetList != 0 {
		log.Fatal("not impl'd")
	}
	if flags&SetAll != 0 {
		log.Fatal("not impl'd")
	}
	return self
}

func (self *Attribute) OnGet(flags AttributeFlags, fn GetEncoderFunc) *Attribute {
	handler := func(req *Request, res Response) { fn(res) }
	if flags&GetSingle != 0 {
		self.OnService(cip.GetAttributeSingle, handler)
	}
	if flags&GetList != 0 {
		self.OnService(cip.GetAttributeList, handler)
	}
	if flags&GetAll != 0 {
		self.OnService(cip.GetAttributesAll, handler)
	}
	return self
}
