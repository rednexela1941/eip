package adapter

import (
	"github.com/rednexela1941/eip/pkg/cip"
)

const MessageRouterObjectRevision cip.UINT = 1

func (self *_Adapter) InitDefaultMessageRouterObject() {
	c := self.AddClass("Message Router", cip.MessageRouterClassCode, MessageRouterObjectRevision)

	c.AddAttribute(1, "ClassRevision", cip.UINTSize).OnGet(
		GetAll|GetSingle,
		func(res Response) { res.Wl(MessageRouterObjectRevision) },
	)
	c.AddAttribute(2, "MaxInstance", cip.UINTSize).OnGet(
		GetSingle,
		func(res Response) { res.Wl(c.HighestInstanceID()) },
	)
	c.AddAttribute(3, "NumInstances", cip.UINTSize).OnGet(
		GetSingle,
		func(res Response) { res.Wl(c.NumberOfInstances()) },
	)
	c.AddAttribute(4, "OptionalAttributeList", cip.UINTSize).OnGet(
		GetAll,
		func(res Response) { res.Wl(cip.UINT(0)) },
	)
	c.AddAttribute(5, "OptionalServiceList", cip.UINTSize).OnGet(
		GetAll,
		func(res Response) { res.Wl(cip.UINT(0)) },
	)
	c.AddAttribute(6, "MaxClassAttributeID", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestAttributeID()) },
	)
	c.AddAttribute(7, "MaxInstanceAttributeID", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestInstanceAttributeID()) },
	)

	c.addDefaultGetAttributesAll()

	i := c.AddInstance(1)

	i.OnService(cip.GetAttributeSingle, func(req *Request, res Response) {
		// TODO: what should go in here?
		res.SetGeneralStatus(cip.StatusAttributeNotSupported)
	})
}
