// Code generated by "stringer -type=LogicalFormatSpecial"; DO NOT EDIT.

package epath

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LogicalFormatSpecialElectronicKey-0]
	_ = x[LogicalFormatSpecialMask-3]
}

const (
	_LogicalFormatSpecial_name_0 = "LogicalFormatSpecialElectronicKey"
	_LogicalFormatSpecial_name_1 = "LogicalFormatSpecialMask"
)

func (i LogicalFormatSpecial) String() string {
	switch {
	case i == 0:
		return _LogicalFormatSpecial_name_0
	case i == 3:
		return _LogicalFormatSpecial_name_1
	default:
		return "LogicalFormatSpecial(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
