package epath

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

/*
### C-1.6
		Encoded Path Compression Rules
		When multiple encoded paths are concatenated the delineation between paths is where a
		segment at a higher level in the hierarchy is encountered.  Multiple encoded paths may be
		compacted when each path shares the same values at the higher levels in the hierarchy.
		Extended Logical Segments shall not be used in compressed paths. When a segment is
		encountered which is at the same or higher level but not at the top level in the hierarchy, the
		preceding higher levels are used for that next encoded path.  The examples below show
		multiple encoded paths in the full and compacted representations.
		Full:
		Class A, Instance A, Attribute A, Class A, Instance A, Attribute B
		Compact:  Class A, Instance A, Attribute A, Attribute B

		Full:
		Class A, Instance A, Attribute A, Class A, Instance B, Attribute A
		Compact:  Class A, Instance A, Attribute A, Instance B, Attribute A
*/

type ApplicationPath struct {
	ClassID         cip.ClassCode
	InstanceID      cip.UINT
	AttributeID     cip.UINT
	MemberID        cip.UINT
	ConnectionPoint cip.UINT
}

func (self *ApplicationPath) Encode(lengthWords cip.USINT) ([]cip.BYTE, error) {
	// TODO: eventually make this good.
	if lengthWords == 0 {
		return nil, fmt.Errorf("cannot Encode 0 length path")
	}

	w := bbuf.New(nil)

	checkLength := func() (bool, error) {
		if lengthWords == 0 {
			return true, nil
		}
		if lengthWords < 0 {
			return true, fmt.Errorf("unable to encode path with length")
		}
		return false, nil
	}

	if self.ClassID < 256 {
		w.Wl(cip.BYTE(0x20))
		w.Wl(cip.USINT(self.ClassID))
		lengthWords -= 1
	} else {
		w.Wl(cip.BYTE(0x21))
		w.Wl(cip.BYTE(0)) // pad
		w.Wl(self.ClassID)
		lengthWords -= 2
	}

	if done, err := checkLength(); done {
		return w.Bytes(), err
	}

	if self.InstanceID < 256 {
		w.Wl(cip.BYTE(0x24))
		w.Wl(cip.USINT(self.InstanceID))
		lengthWords -= 1
	} else {
		w.Wl(cip.BYTE(0x25))
		w.Wl(cip.BYTE(0)) // pad
		w.Wl(self.InstanceID)
		lengthWords -= 2
	}
	if done, err := checkLength(); done {
		return w.Bytes(), err
	}
	return nil, fmt.Errorf("incomplete")
}

func (self *ApplicationPath) GetInstanceIDOrConnectionPoint() cip.UINT {
	if self.InstanceID == 0 && self.ConnectionPoint != 0 {
		return self.ConnectionPoint
	}
	return self.InstanceID
}

func (self *ApplicationPath) DebugString() string {
	s := ""
	if self.ClassID != 0 {
		s += self.ClassID.String()
		i := self.InstanceID
		if i == 0 && self.ConnectionPoint != 0 {
			s += fmt.Sprintf("ConnPt%d", self.ConnectionPoint)
		} else {
			s += fmt.Sprintf(".Instance%d", i)
		}
	}
	if self.AttributeID != 0 {
		s += fmt.Sprintf(".Attr%d", self.AttributeID)
	}
	if self.MemberID != 0 {
		s += fmt.Sprintf(".Member%d", self.MemberID)
	}
	return s
}

func (self *ApplicationPath) setValue(t LogicalSegType, value cip.UINT) error {
	switch t {
	case LogicalSegTypeClassID:
		self.ClassID = cip.ClassCode(value)
	case LogicalSegTypeInstanceID:
		self.InstanceID = value
	case LogicalSegTypeAttributeID:
		self.AttributeID = value
	case LogicalSegTypeConnectionPoint:
		self.ConnectionPoint = value
	case LogicalSegTypeMemberID:
		self.MemberID = value
	default:
		return fmt.Errorf("unhandled logical type %s", t.String())
	}
	return nil
}

type appPathParseLevel int

const (
	appPathParseLevelStart appPathParseLevel = iota
	appPathParseLevelClass
	appPathParseLevelInstance
	appPathParseLevelAttribute
	appPathParseLevelMember
)

type appPathItem struct {
	l LogicalSegType
	v cip.UINT
}

func (self LogicalSegType) toLevel() appPathParseLevel {
	switch self {
	case LogicalSegTypeClassID:
		return appPathParseLevelClass
	case LogicalSegTypeInstanceID, LogicalSegTypeConnectionPoint:
		return appPathParseLevelInstance
	case LogicalSegTypeMemberID:
		return appPathParseLevelMember
	default:
		return appPathParseLevelAttribute
	}
}

type appPathStack struct {
	items []appPathItem
}

func (self *appPathStack) len() int {
	return len(self.items)
}

func (self *appPathStack) pop() {
	lastIndex := len(self.items) - 1
	if lastIndex < 0 {
		return
	}
	self.items = self.items[:lastIndex]
}

func (self *appPathStack) popUntil(level appPathParseLevel) {
	for {
		if len(self.items) == 0 {
			return
		}
		if self.items[len(self.items)-1].l.toLevel() < level {
			return
		}
		self.pop()
	}
}

func (self *appPathStack) push(l LogicalSegType, v cip.UINT) {
	self.items = append(self.items, appPathItem{l: l, v: v})
}

func (self *appPathStack) toPath() (ApplicationPath, error) {
	ap := ApplicationPath{}
	for _, i := range self.items {
		if err := ap.setValue(i.l, i.v); err != nil {
			return ap, err
		}
	}
	return ap, nil
}

func parseApplicationPaths(r *bbuf.Buffer, padded bool) ([]ApplicationPath, error) {
	aps := make([]ApplicationPath, 0, 1)
	currentLevel := appPathParseLevelStart
	stack := appPathStack{
		items: make([]appPathItem, 0),
	}
	for r.Len() > 0 {
		lseg, value, err := parseLogicalSegment(r, padded)
		if err != nil {
			return nil, err
		}
		level := lseg.toLevel()
		if level <= currentLevel {
			add, err := stack.toPath()
			if err != nil {
				return nil, err
			}
			aps = append(aps, add)
			stack.popUntil(level)
		}
		currentLevel = level
		stack.push(lseg, value)
		// shift on and off properties until we get to current parse level.
	}
	if stack.len() > 0 {
		add, err := stack.toPath()
		if err != nil {
			return nil, err
		}
		aps = append(aps, add)
	}
	return aps, nil
}
