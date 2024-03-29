// Code generated by "stringer -type=ItemID"; DO NOT EDIT.

package cpf

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AddressNull-0]
	_ = x[CIPIdentity-12]
	_ = x[SecurityInfo-134]
	_ = x[EtherNetIPCapability-135]
	_ = x[EtherNetIPUsage-136]
	_ = x[AddressConnected-161]
	_ = x[ConnectedTransportPacket-177]
	_ = x[UnconnectedMessage-178]
	_ = x[ListServicesResponse-256]
	_ = x[SockaddrInfoOtoT-32768]
	_ = x[SockaddrInfoTtoO-32769]
	_ = x[SequencedAddress-32770]
	_ = x[UnconnectedMessageOverUDP-32771]
}

const (
	_ItemID_name_0 = "AddressNull"
	_ItemID_name_1 = "CIPIdentity"
	_ItemID_name_2 = "SecurityInfoEtherNetIPCapabilityEtherNetIPUsage"
	_ItemID_name_3 = "AddressConnected"
	_ItemID_name_4 = "ConnectedTransportPacketUnconnectedMessage"
	_ItemID_name_5 = "ListServicesResponse"
	_ItemID_name_6 = "SockaddrInfoOtoTSockaddrInfoTtoOSequencedAddressUnconnectedMessageOverUDP"
)

var (
	_ItemID_index_2 = [...]uint8{0, 12, 32, 47}
	_ItemID_index_4 = [...]uint8{0, 24, 42}
	_ItemID_index_6 = [...]uint8{0, 16, 32, 48, 73}
)

func (i ItemID) String() string {
	switch {
	case i == 0:
		return _ItemID_name_0
	case i == 12:
		return _ItemID_name_1
	case 134 <= i && i <= 136:
		i -= 134
		return _ItemID_name_2[_ItemID_index_2[i]:_ItemID_index_2[i+1]]
	case i == 161:
		return _ItemID_name_3
	case 177 <= i && i <= 178:
		i -= 177
		return _ItemID_name_4[_ItemID_index_4[i]:_ItemID_index_4[i+1]]
	case i == 256:
		return _ItemID_name_5
	case 32768 <= i && i <= 32771:
		i -= 32768
		return _ItemID_name_6[_ItemID_index_6[i]:_ItemID_index_6[i+1]]
	default:
		return "ItemID(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
