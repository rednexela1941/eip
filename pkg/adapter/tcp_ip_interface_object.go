package adapter

import (
	"time"

	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/network"
)

// See TCP IP interface object attribute 5
type InterfaceConfiguration struct {
	network.InterfaceAddr
	Gateway                network.IPv4
	NameServer             network.IPv4
	NameServer2            network.IPv4
	DomainNameLengthAndPad cip.UINT // zero, because we don't support domain names yet.
}

const TCPIPDefaultEncapsulationInvactivityTimeout cip.UINT = 120 // seconds.

// See Volume 2: 5-4 "TCP/IP Interface Object"
const TCPIPInterfaceObjectRevision cip.UINT = 4

// back to page 92.
// See Volume 2: 5-4.3.2.1
// And Volume 2: Table 5-4.4
type TCPIPInterfaceStatus cip.DWORD

const (
	TCPIPNotConfigured TCPIPInterfaceStatus = 0b000
	TCPIPDHCP          TCPIPInterfaceStatus = 0b001
	TCPIPStatic        TCPIPInterfaceStatus = 0b010

	TCPIPMcastPending              TCPIPInterfaceStatus = 1 << 4
	TCPIPInterfaceConfigPending    TCPIPInterfaceStatus = 1 << 5
	TCPIPAcdStatusConflictDetected TCPIPInterfaceStatus = 1 << 6
	TCPIPAcdFaultDetected          TCPIPInterfaceStatus = 1 << 7
	TCPIPIANAPortChangePending     TCPIPInterfaceStatus = 1 << 8
	TCPIPIANAProtocolChangePending TCPIPInterfaceStatus = 1 << 9
)

// See Volume 2: 5-4.3.2.2 "Configuration Capability"
type TCPIPConfigurationCapability cip.DWORD

const (
	TCPIPConfigBOOTP                  TCPIPConfigurationCapability = 1 << 0
	TCPIPConfigDNS                    TCPIPConfigurationCapability = 1 << 1
	TCPIPConfigDHCP                   TCPIPConfigurationCapability = 1 << 2
	TCPIPConfigDHCPDNSUpdate          TCPIPConfigurationCapability = 0 << 3 // shall be zero
	TCPIPConfigSetable                TCPIPConfigurationCapability = 1 << 4 // can set
	TCPIPConfigHardwareConfigurable   TCPIPConfigurationCapability = 1 << 5
	TCPIPInterfaceConfigRequiresReset TCPIPConfigurationCapability = 1 << 6
	TCPIPConfigAcdCapable             TCPIPConfigurationCapability = 1 << 7
)

// See Volume 2: 5-4.3.2.3 "Configuration Control"
type TCPIPConfigurationControl cip.DWORD

const (
	TCPIPControlStatic TCPIPConfigurationControl = 0
	TCPIPControlBOOTP  TCPIPConfigurationControl = 1
	TCPIPControlDHCP   TCPIPConfigurationControl = 2

	TCPIPControlDNSEnable TCPIPConfigurationControl = 1 << 4
)

// See Volume 2: 5-4.3.2 "Instance Attributes"
func (self *Adapter) GetEncapsulationInactivityTimeout() time.Duration {
	// 120 seconds is the default timeout
	return 120 * time.Second
}

func (self *Adapter) GetTCPIPInterfaceClass() *Class {
	c, ok := self.Classes[cip.TCPIPInterfaceClassCode]
	if ok {
		return c
	}

	c = self.AddClass(
		"TCP/IP Interface",
		cip.TCPIPInterfaceClassCode,
		TCPIPInterfaceObjectRevision,
	)

	c.AddAttribute(1, "ClassRevision", cip.UINTSize).OnGet(
		GetSingle|GetAll,
		func(res Response) { res.Wl(TCPIPInterfaceObjectRevision) },
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

	c.addDefaultGetAttributesAll()

	return c
}
