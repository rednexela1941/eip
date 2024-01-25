package adapter

import "github.com/rednexela1941/eip/pkg/cip"

const EthernetLinkObjectRevision cip.UINT = 4

// See Volume 2: 5-5.3.2.2 "Interface Flags - Attribute 2"
type EthernetLinkInterfaceFlags cip.DWORD

const (
	EthernetLinkFlagLinkInactive                   EthernetLinkInterfaceFlags = 0
	EthernetLinkFlagLinkActive                     EthernetLinkInterfaceFlags = 1
	EthernetLinkFlagHalfDuplex                     EthernetLinkInterfaceFlags = 0 << 1
	EthernetLinkFlagFullDuplex                     EthernetLinkInterfaceFlags = 1 << 1
	EthernetLinkFlagAutoNegotiationInProgress      EthernetLinkInterfaceFlags = 0 << 2
	EthernetLinkFlagAutoNegotiationFailed          EthernetLinkInterfaceFlags = 1 << 2
	EthernetLinkFlagAutoNegotiationFailedWithSpeed EthernetLinkInterfaceFlags = 2 << 2
	EthernetLinkFlagAutoNegotiationSuccess         EthernetLinkInterfaceFlags = 3 << 2
	EthernetLinkFlagAutoNegotitationNotAttempted   EthernetLinkInterfaceFlags = 4 << 2
	EthernetLinkFlagManualSettingRequiresReset     EthernetLinkInterfaceFlags = 1 << 5
	EthernetLinkFlagLocalHardwareFault             EthernetLinkInterfaceFlags = 1 << 6
)

// See Volume 2: 5-5.3.2.7
type EthernetLinkType cip.USINT

const (
	EthernetLinkTypeUnknown         EthernetLinkType = 0
	EthernetLinkTypeInternal        EthernetLinkType = 1
	EthernetLinkTypeTwistedPair     EthernetLinkType = 2
	EthernetLinkTypeOpticalFiber    EthernetLinkType = 3
	EthernetLinkTypeInCabinetRibbon EthernetLinkType = 4
)

// See Volume 2: 5-5.3.2.8
type EthernetLinkState cip.USINT

const (
	EthernetLinkStateUnknown  EthernetLinkState = 0
	EthernetLinkStateReader   EthernetLinkState = 1
	EthernetLinkStateDisabled EthernetLinkState = 2
	EthernetLinkStateTesting  EthernetLinkState = 3
)

// See Volume 2: 5-5.3.2.11 "Interface Capability - Attribute 11"
type EthernetLinkCapability cip.DWORD

const (
	EthernetLinkCapabilityManualSettingRequiresReset EthernetLinkCapability = 1
	EthernetLinkCapabilityAutoNegotiate              EthernetLinkCapability = 1 << 0
	EthernetLinkCapabilityAutoMDIX                   EthernetLinkCapability = 1 << 2
	EthernetLinkCapabilityManualSpeedDuplex          EthernetLinkCapability = 1 << 3
)

// See Volume 2: 5-5.3.2 (Attribute 11, inside the struct)
type InterfaceDuplexMode cip.USINT

const (
	InterfaceDuplexModeHalf InterfaceDuplexMode = 0
	InterfaceDuplexModeFull InterfaceDuplexMode = 1
)

// Return
func (self *Adapter) GetEthernetLinkClass() *Class {
	c, ok := self.Classes[cip.EthernetLinkClassCode]
	if ok {
		return c
	}

	c = self.AddClass("Ethernet Link", cip.EthernetLinkClassCode, EthernetLinkObjectRevision)

	c.AddAttribute(1, "ClassRevision", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(EthernetLinkObjectRevision) },
	)
	c.AddAttribute(2, "MaxInstance", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(c.HighestInstanceID()) },
	)
	c.AddAttribute(3, "NumInstances", cip.UINTSize).OnGet(
		GetSingle|GetAll,
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
