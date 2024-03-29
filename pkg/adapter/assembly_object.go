package adapter

import (
	"fmt"
	"log"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/param"
)

type AssemblyInstance struct {
	*Instance
	Name       string // for EDS
	Parameters []*param.AssemblyParam
}

const AssemblyObjectRevision cip.UINT = 3

// See Volume 1: 5A-5.3 "Connection Points"
// In particular:
// Connection Points of the Assembly Object are identical to the Data Attribute (#3) of the
// Instances.  For example, Connection Point 4 of the Assembly Object is the same as Instance 4,
// Attribute #3.  Specifying a path of “20 04 24 xx 30 03” is the same as “20 04 2C xx”.

// GetAssemblyClass get the underlying assembly class, or create default if it doesn't exist. Shouldn't have to call this, just use AddAssemblyInstance.
func (self *_Adapter) GetAssemblyClass() *Class {
	c, ok := self.Classes[cip.AssemblyClassCode]
	if ok {
		return c
	}

	// Add default assembly object.
	c = self.AddClass("Assembly", cip.AssemblyClassCode, AssemblyObjectRevision)

	c.AddAttribute(1, "ClassRevision", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(AssemblyObjectRevision) },
	)
	c.AddAttribute(2, "MaxInstance", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) {
			res.Wl(c.HighestInstanceID())
		},
	)
	c.AddAttribute(3, "NumInstances", cip.UINTSize).OnGet(
		GetSingle,
		func(res Response) { res.Wl(c.NumberOfInstances()) },
	)
	c.AddAttribute(6, "MaxClassAttributeID", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestAttributeID()) },
	)
	c.AddAttribute(7, "MaxInstanceAttributeID", cip.UDINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestInstanceAttributeID()) },
	)

	c.OnService(cip.GetAttributesAll, func(req *Request, res Response) {
		res.SetGeneralStatus(cip.StatusServiceNotSupported)
	})

	return c
}

func (self *AssemblyInstance) handleSet(req *Request, res Response) {
	_, ok := self.Attributes[3]
	if !ok {
		res.SetGeneralStatus(cip.StatusAttributeNotSupported)
		return
	}

	size := self.GetSize()

	if len(req.Request.RequestData) < int(size) {
		res.SetGeneralStatus(cip.StatusNotEnoughData)
		return
	}
	if len(req.Request.RequestData) > int(size) {
		res.SetGeneralStatus(cip.StatusTooMuchData)
		return
	}

	if err := self.ReadFrom(bbuf.New(req.Request.RequestData)); err != nil {
		log.Println(err)
		res.SetGeneralStatus(cip.StatusAttributeNotSetable)
		return
	}
}

func (self *AssemblyInstance) handleGet(res Response) {
	_, ok := self.Attributes[3]
	if !ok {
		res.SetGeneralStatus(cip.StatusAttributeNotSupported)
		return
	}

	if err := self.WriteTo(res); err != nil {
		res.SetGeneralStatus(cip.StatusAttributeNotGettable)
		return
	}
}

func (self *AssemblyInstance) ReadFrom(r bbuf.Reader) error {
	for _, p := range self.Parameters {
		startLen := r.Len()
		if err := p.ReadFrom(r); err != nil {
			return err
		}
		if r.Error() != nil {
			return r.Error()
		}
		endLen := r.Len()
		diff := startLen - endLen
		if diff != int(p.DataType.Size) {
			return fmt.Errorf(
				"read wrong amount of bytes for param %s (read %d need %d)",
				diff,
				p.DataType.Size,
			)
		}
	}
	return nil
}

func (self *AssemblyInstance) WriteTo(w bbuf.Writer) error {
	for _, p := range self.Parameters {
		startLen := w.Len()
		if err := p.WriteTo(w); err != nil {
			return err
		}
		if w.Error() != nil {
			return w.Error()
		}
		endLen := w.Len()
		diff := endLen - startLen
		if diff != int(p.DataType.Size) {
			return fmt.Errorf(
				"wrote wrong amount of bytes for param %s (need %d got %d)",
				p.Name, p.DataType.Size, diff,
			)
		}
	}
	return nil
}

func (self *_Adapter) AddAssemblyInstance(name string, instanceID cip.UINT) *AssemblyInstance {
	c := self.GetAssemblyClass()
	i := c.AddInstance(instanceID)

	ai := new(AssemblyInstance)
	ai.Instance = i
	ai.Parameters = make([]*param.AssemblyParam, 0)
	ai.Name = name

	self.AssemblyInstances = append(self.AssemblyInstances, ai)

	i.AddAttribute(3, "Data", 0 /* temp size */).OnGet(
		GetSingle,
		ai.handleGet,
	).OnSet(
		SetSingle,
		ai.handleSet,
	)
	i.AddAttribute(4, "Size", cip.UINTSize).OnGet(
		GetSingle,
		func(res Response) { res.Wl(ai.GetSize()) },
	)

	return ai
}

func (self *AssemblyInstance) GetSize() cip.UINT {
	size := cip.UINT(0)
	for _, p := range self.Parameters {
		size += p.DataType.Size
	}
	return size
}

// Check that types are aligned on proper boundaries.
// TODO: link to rockwell document on this.
func (self *AssemblyInstance) CheckParamAlignment() error {
	offset := 0
	for _, p := range self.Parameters {
		size := int(p.DataType.Size)
		mod := offset % size
		if mod != 0 {
			return fmt.Errorf("%s is not aligned (to fix add %d pad bytes)", p.Name, mod)
		}
		offset += size
	}
	return nil
}

func (self *AssemblyInstance) AddParam(p *param.AssemblyParam) {
	self.Parameters = append(self.Parameters, p)
	if p.DataType.Size == 1 {
		// no need to check this here.
		return
	}
	// this loop is a bit repetitice.
	if err := self.CheckParamAlignment(); err != nil {
		log.Fatal(err)
	}
}

// Add a pad byte parameter
func (self *AssemblyInstance) AddPadByteParam() {
	self.AddSINTParam("Padding").OnGet(
		func() cip.SINT { return 0 },
	).OnSet(
		func(v cip.SINT) { return },
	)
}
