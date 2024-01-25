package adapter

import (
	"log"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

const IdentityObjectRevision cip.UINT = 2

// silence warning.
func noop(a interface{}) { return }

const (
	ResetTypePowerCycle                 cip.USINT = 0
	ResetTypeFactoryDefaults            cip.USINT = 1
	ResetTypeFactoryDefaultsExceptComms cip.USINT = 2
)

func (self *_Adapter) InitDefaultIdentityObject() {
	c := self.AddClass("Identity", cip.IdentityClassCode, IdentityObjectRevision)

	c.AddAttribute(1, "ClassRevision", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) {
			res.Wl(IdentityObjectRevision)
		},
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

	c.addDefaultGetAttributesAll()
	// c.addDefaultGetAttributesList()

	i := c.AddInstance(1) // add first identity object instance.
	i.addDefaultGetAttributesList()
	i.addDefaultGetAttributesAll()
	// i.addDefaultSetAttributeList()

	i.AddAttribute(1, "VendorID", cip.UINTSize).OnGet(
		GetFull,
		func(res Response) { res.Wl(self.Identity.VendorID) },
	)
	i.AddAttribute(2, "DeviceType", cip.UINTSize).OnGet(
		GetFull,
		func(res Response) { res.Wl(self.Identity.DeviceType) },
	)
	i.AddAttribute(3, "ProductCode", cip.UINTSize).OnGet(
		GetFull,
		func(res Response) { res.Wl(self.Identity.ProductCode) },
	)
	i.AddAttribute(4, "DeviceRevision", cip.USINTSize+cip.USINTSize).OnGet(
		GetFull,
		func(res Response) { res.Wl(self.Identity.Revision) },
	)
	i.AddAttribute(5, "DeviceStatus", cip.WORDSize).OnGet(
		GetFull,
		func(res Response) { res.Wl(self.Identity.Status) },
	)
	i.AddAttribute(6, "SerialNumber", cip.UDINTSize).OnGet(
		GetFull,
		func(res Response) { res.Wl(self.Identity.SerialNumber) },
	)
	i.AddAttribute(7, "ProductName", 1+len(self.Identity.ProductName)).OnGet(
		GetFull,
		func(res Response) {
			// encode name as SHORT_STRING
			if err := bbuf.WShortString(res, self.Identity.ProductName); err != nil {
				log.Fatal(err)
			}
		},
	)

	resetHandler := func(req *Request, res Response) {
		resetType := cip.USINT(0)
		if len(req.RequestData) > 0 {
			resetType = req.RequestData[0]
		}
		switch resetType {
		case ResetTypePowerCycle:
			// Power Cycle.
			res.SetGeneralStatus(cip.StatusDeviceStateConflict)
			return
		case ResetTypeFactoryDefaults:
			// Return to Factory Defaults.
			res.SetGeneralStatus(cip.StatusDeviceStateConflict)
		case ResetTypeFactoryDefaultsExceptComms:
			// Return to Factory Defaults except communication parameters.
			res.SetGeneralStatus(cip.StatusDeviceStateConflict)
		default:
			res.SetGeneralStatus(cip.StatusInvalidParameter)
			return
		}
	}

	// c.OnService(cip.Reset, resetHandler)
	i.OnService(cip.Reset, resetHandler)
}
