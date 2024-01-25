// Code generated by "stringer -type=Command"; DO NOT EDIT.

package encap

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NOP-0]
	_ = x[ListServices-4]
	_ = x[ListIdentity-99]
	_ = x[ListInterfaces-100]
	_ = x[RegisterSession-101]
	_ = x[UnRegisterSession-102]
	_ = x[SendRRData-111]
	_ = x[SendUnitData-112]
	_ = x[StartDTLS-200]
}

const (
	_Command_name_0 = "NOP"
	_Command_name_1 = "ListServices"
	_Command_name_2 = "ListIdentityListInterfacesRegisterSessionUnRegisterSession"
	_Command_name_3 = "SendRRDataSendUnitData"
	_Command_name_4 = "StartDTLS"
)

var (
	_Command_index_2 = [...]uint8{0, 12, 26, 41, 58}
	_Command_index_3 = [...]uint8{0, 10, 22}
)

func (i Command) String() string {
	switch {
	case i == 0:
		return _Command_name_0
	case i == 4:
		return _Command_name_1
	case 99 <= i && i <= 102:
		i -= 99
		return _Command_name_2[_Command_index_2[i]:_Command_index_2[i+1]]
	case 111 <= i && i <= 112:
		i -= 111
		return _Command_name_3[_Command_index_3[i]:_Command_index_3[i+1]]
	case i == 200:
		return _Command_name_4
	default:
		return "Command(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
