package linux

import (
	"log"
	"net"
	"os"

	"github.com/rednexela1941/eip/pkg/adapter"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/network"
)

type NetworkInterface struct {
	Interface *net.Interface
}

func CreateInterface(name string) (*NetworkInterface, error) {
	i := new(NetworkInterface)
	compilerCheck(i)

	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	i.Interface = iface
	return i, nil
}

func (self *NetworkInterface) GetInterface() *net.Interface {
	return self.Interface
}

func (self *NetworkInterface) GetSpeed() adapter.Mbps {
	return 100
}

func (self *NetworkInterface) GetFlags() adapter.EthernetLinkInterfaceFlags {
	flags := adapter.EthernetLinkFlagLinkActive
	flags |= adapter.EthernetLinkFlagFullDuplex
	flags |= adapter.EthernetLinkFlagAutoNegotiationSuccess
	// flags |= adapter.EthernetLinkFlagManualSettingRequiresReset
	return flags
}

func (self *NetworkInterface) GetType() adapter.EthernetLinkType {
	return adapter.EthernetLinkTypeTwistedPair
}

func (self *NetworkInterface) GetCapabilities() adapter.InterfaceCapabilities {
	speeds := make([]adapter.SpeedDuplex, 0)
	sp := 10
	for sp <= 1000 {
		speeds = append(speeds, adapter.SpeedDuplex{
			Speed:      cip.UINT(sp),
			DuplexMode: adapter.InterfaceDuplexModeHalf,
		})
		speeds = append(speeds, adapter.SpeedDuplex{
			Speed:      cip.UINT(sp),
			DuplexMode: adapter.InterfaceDuplexModeFull,
		})
		sp *= 10
	}
	return adapter.InterfaceCapabilities{
		CapabilityBits: adapter.EthernetLinkCapabilityAutoNegotiate |
			adapter.EthernetLinkCapabilityAutoMDIX,
		SpeedDuplexArray: speeds,
	}
}

func (self *NetworkInterface) GetStatus() adapter.TCPIPInterfaceStatus {
	return adapter.TCPIPStatic
}

func (self *NetworkInterface) GetTCPConfigCapability() adapter.TCPIPConfigurationCapability {
	return 0
}

func (self *NetworkInterface) GetTCPConfigControl() adapter.TCPIPConfigurationControl {
	return adapter.TCPIPControlStatic
}

func (self *NetworkInterface) GetHostname() string {
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func (self *NetworkInterface) GetAddresses() []*network.InterfaceAddr {
	addrs, err := self.Interface.Addrs()
	if err != nil {
		log.Fatal(err)
	}
	result := make([]*network.InterfaceAddr, 0)
	for _, ra := range addrs {
		a, ok := ra.(*net.IPNet)
		if !ok {
			continue
		}
		ia, err := network.FromIPNet(a)
		if err != nil {
			continue
		}
		result = append(result, ia)
	}
	return result
}

func compilerCheck(i *NetworkInterface) adapter.NetworkInterface {
	return i
}
