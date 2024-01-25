package epath

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

// see section C-1.4 in volume 1

type Segment cip.USINT

//go:generate stringer -type=SegmentType
type SegmentType cip.USINT

/***
When an extended logical segment is included within a Padded Path and the Logical Format is
8 bit, a pad byte shall be added after the Logical Value (the 16-bit and 32-bit formats are
identical to the Packed Path) and shall be set to 0.  For all other logical segments, when
included within a Padded Path, the 16-bit and 32-bit logical formats shall have a pad inserted
between the segment type byte and the Logical Value (the 8-bit format is identical to the
Packed Path). The pad byte shall be set to zero.
***/

const (
	SegmentTypePort                SegmentType = 0b000 << 5
	SegmentTypeLogical             SegmentType = 0b001 << 5
	SegmentTypeNetwork             SegmentType = 0b010 << 5
	SegmentTypeSymbolic            SegmentType = 0b011 << 5
	SegmentTypeData                SegmentType = 0b100 << 5
	SegmentTypeDataTypeConstructed SegmentType = 0b101 << 5
	SegmentTypeDataTypeElementary  SegmentType = 0b110 << 5
	SegmentTypeReserved            SegmentType = 0b111 << 5

	SegmentTypeMask SegmentType = 0b111 << 5
)

type PortSegment cip.USINT

const (
	// If this bit is high, then the next byte is the port segement size in byte.
	// if it is not set, the link address is 1 byte, which follows.
	// see Vol 1 C-1.4.1
	PortSegmentExtendedLinkAddressSizeMask PortSegment = 0b1 << 4
	PortSegmentPortIdentifierMask          PortSegment = 0b1111
	// luckily, nobody cares about port segments.
)

//go:generate stringer -type=LogicalSegType
type LogicalSegType cip.USINT

const (
	LogicalSegTypeClassID         LogicalSegType = 0b000 << 2
	LogicalSegTypeInstanceID      LogicalSegType = 0b001 << 2
	LogicalSegTypeMemberID        LogicalSegType = 0b010 << 2
	LogicalSegTypeConnectionPoint LogicalSegType = 0b011 << 2
	LogicalSegTypeAttributeID     LogicalSegType = 0b100 << 2
	LogicalSegTypeSpecial         LogicalSegType = 0b101 << 2
	LogicalSegTypeServiceID       LogicalSegType = 0b110 << 2
	LogicalSegTypeExtendedLogical LogicalSegType = 0b111 << 2 // see Vol1 Table C-1.3

	LogicalSegTypeMask LogicalSegType = 0b111 << 2
)

//go:generate stringer -type=LogicalFormat
type LogicalFormat cip.USINT

const (
	LogicalFormat8Bit     LogicalFormat = 0b00
	LogicalFormat16Bit    LogicalFormat = 0b01
	LogicalFormat32Bit    LogicalFormat = 0b10 // only allowed for certain types.
	LogicalFormatReserved LogicalFormat = 0b11

	LogicalFormatMask LogicalFormat = 0b11
)

//go:generate stringer -type=LogicalFormatSpecial
type LogicalFormatSpecial cip.USINT

const (
	LogicalFormatSpecialElectronicKey LogicalFormatSpecial = 0b00
	LogicalFormatSpecialMask          LogicalFormatSpecial = 0b11
)

//go:generate stringer -type=LogicalFormatServiceID
type LogicalFormatServiceID cip.USINT

const (
	LogicalFormatServiceID8Bit LogicalFormatServiceID = 0b00
	LogicalFormatServiceIDMask LogicalFormatServiceID = 0b11
)

func (self Segment) IsElectronicKey() bool {
	// alternatively, return 0x34 == self.
	t := self.Type()
	if t != SegmentTypeLogical {
		return false
	}
	lt := self.LogicalType()
	if lt != LogicalSegTypeSpecial {
		return false
	}
	f := LogicalFormatSpecial(self.LogicalFormat())
	return f == LogicalFormatSpecialElectronicKey
}

func (self Segment) Type() SegmentType {
	return SegmentType(self) & SegmentTypeMask
}

func (self Segment) LogicalType() LogicalSegType {
	return LogicalSegType(self) & LogicalSegTypeMask
}

func (self Segment) LogicalFormat() LogicalFormat {
	return LogicalFormat(self) & LogicalFormatMask
}

/***
	C-1.5
		Segment Definition Hierarchy
		In general, the definition of any rules related to the use of segments is defined within the Object
		Class and/or Device Profile. When symbolic and/or logical segments are used to construct an
		application path the Backus-Naur Form definition is:
		application_path ::= CHOICE {symbolic_application_path, class_application_path,
		assembly_class_application_path}
		symbolic_application_path ::= symbolic_segment [Connection_Point] [member_specification]
		[bit_specification]
		symbolic_segment ::= CHOICE {Symbolic_Segment, ANSI_Extended_Symbol_Segment}
		assembly_class_application_path ::= 20 04 [assembly_attribute_specification]
		class_application_path ::= Class_ID1 [item_specification]
		item_specification ::= CHOICE {attribute_specification, connection_point_specification}
		assembly_attribute_specification ::= CHOICE {Instance ID, Connection Point} [[Attribute ID]
		[member_specification] [bit_specification]]
		attribute_specification ::= Instance ID [[Attribute ID] [member_specification]
		[bit_specification]]
		connection_point_specification ::= [Instance ID] Connection Point [member_specification]
		[bit_specification]
		member_specification ::= CHOICE {Member_ID, extended_member_specification}
		extended_member_specification ::= SEQUENCE of {Array_Index,
		indirect_array_specification, Structure_Member_Number, Structure_Member_Handle}
		indirect_array_specification ::= Indirect_Array_Index  application_path
		bit_specification ::= CHOICE {Bit_Index, indirect_index_bit_specification}
		indirect_bit_specification ::= Indirect_Bit_Index  application_path

		NOTE: A Symbolic_application_path must resolve, at run time, to a Class_application_path.
		NOTE: An Indirect_Array_Index and Indirect_Bit_Index application_path must evaluate to a
		positive integer.
		The depth within the application path need only proceed to the degree required by its
		application.
		The following are examples of valid application paths.
		• Class ID,
		• Class ID, Instance ID
		• Class ID, Instance ID, Attribute ID
		• Class ID, Instance ID, Attribute ID, Member ID
		• Class ID, Connection Point
		• Class ID, Connection Point, Member ID
		• Class ID, Instance ID, Connection Point2
		• Class ID, Instance ID, Connection Point, Member ID 2
		• Symbolic ID
		• Symbolic ID, Member ID
		• Symbolic ID, Connection Point 3
		• Symbolic ID, Connection Point, Member ID 3
		• Symbolic ID, Array Index, Structure Member Number, Array Index, Bit Index
		• Symbolic ID, Indirect_Array_Index, [Class_ID, Instance_ID, Attribute_ID],
		Structure_Member_Number, Array_Index, Bit_Index
***/

func parseLogicalSegment(r bbuf.Reader, padded bool) (LogicalSegType, cip.UINT, error) {
	// Only application paths for right now.
	var segment Segment
	r.Rl(&segment)

	if segment.Type() != SegmentTypeLogical {
		return 0, 0, fmt.Errorf("unhandled segment type %s", segment.Type().String())
	}

	logical := segment.LogicalType()
	format := segment.LogicalFormat()

	value := cip.UINT(0)
	if padded {
		if format == LogicalFormat16Bit || format == LogicalFormat32Bit {
			pad := byte(0)
			r.Rl(&pad)
		}
	}

	switch format {
	case LogicalFormat8Bit:
		bval := byte(0)
		r.Rl(&bval)
		value = cip.UINT(bval)
	case LogicalFormat16Bit:
		r.Rl(&value)
	case LogicalFormat32Bit:
		// only allowed for InstanceID, MemberID, ConnectionPoint, Array Index, Bit Index, Structure Memeber Number Handle types.
		// not supported right now
		return 0, 0, fmt.Errorf("invalid logical format %s", format.String())
	default:
		return 0, 0, fmt.Errorf("invalid logical format %s", format.String())
	}

	return logical, value, r.Error()
}
