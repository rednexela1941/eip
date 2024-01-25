package adapter

import (
	"fmt"
	"log"

	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cm"
	"github.com/rednexela1941/eip/pkg/epath"
)

//go:generate stringer -type=ConnectionPointType
//go:generate stringer -type=RealTimeFormat

type (
	ConnectionPointType int
	RealTimeFormat      int

	ConnectionPoint struct {
		Name       string // for EDS
		Transport  cm.TransportClassAndTrigger
		Type       ConnectionPointType
		Config     *AssemblyInstance
		Input      *AssemblyInstance
		Output     *AssemblyInstance
		OtoTFormat RealTimeFormat
		TtoOFormat RealTimeFormat
	}
)

const (
	// See Volume 1: 3-6.5.1
	ListenOnly ConnectionPointType = 1 << 0
	// See Volume 1: 3-5.6.2
	InputOnly ConnectionPointType = 1 << 1
	// See Volume 1: 3-5.6.3
	ExclusiveOwner ConnectionPointType = 1 << 2
	// See Volume 1: 3-5.6.4
	RedundantOwner ConnectionPointType = 1 << 3
)

const (
	// Volume 1: 3-6.1 "Real time formats including RUN/IDLE notification"
	// Volume 1: 3-6.1.1: Modeless Format
	ModelessFormat RealTimeFormat = 1 << 0
	// Volume 1: 3-6.1.2: Zero Length Data Format
	ZeroLengthFormat RealTimeFormat = 1 << 1
	// Volume 1: 3-6.1.3: Heartbeat Format
	HeartbeatFormat RealTimeFormat = 1 << 2
	// Volume 1: 3-6.1.4: 32-Bit Header Format
	Header32BitFormat RealTimeFormat = 1 << 3
)

func (self RealTimeFormat) Size(class cm.TransportClassAndTrigger) cip.UINT {
	c := class.TransportClass()

	switch self {
	case ModelessFormat, ZeroLengthFormat, HeartbeatFormat:
		if c == cm.Class0 {
			return 0
		}
		if c == cm.Class1 {
			return 2 // sequence count
		}
	case Header32BitFormat:
		if c == cm.Class0 {
			return 4
		}
		if c == cm.Class1 {
			return 4 + 2 // header + sequence count.
		}
	default:
		log.Fatalf("unknown format %s", self.String())
	}
	log.Fatalf("invalid transport class %d", c)
	return 0
}

// AddConnectionPoint: instances may be nil, depending on type.
// for Class 0 / Class 1 connections.
func (self *_Adapter) AddConnectionPoint(
	name string,
	transport cm.TransportClassAndTrigger,
	connectionPointType ConnectionPointType,
	config *AssemblyInstance,
	input *AssemblyInstance,
	output *AssemblyInstance,
	o2tFormat RealTimeFormat,
	t2OFormat RealTimeFormat,
) {

	self.ConnectionPoints = append(self.ConnectionPoints, ConnectionPoint{
		Name:       name,
		Transport:  transport,
		Type:       connectionPointType,
		Input:      input,
		Config:     config,
		Output:     output,
		OtoTFormat: o2tFormat,
		TtoOFormat: t2OFormat,
	})
}

// Only use for non-null non-matching.
func (self *_Adapter) GetMatchingConnectionPoint(
	freq *cm.SharedForwardOpenRequest,
) (*ConnectionPoint, error) {
	config, ok := freq.GetConfigPath()
	if !ok {
		return nil, fmt.Errorf("invalid config path")
	}
	prod, ok := freq.GetProductionPath()
	if !ok {
		return nil, fmt.Errorf("invalid production path")
	}
	cons, ok := freq.GetConsumptionPath()
	if !ok {
		return nil, fmt.Errorf("invalid consumption path")
	}

	doesMatch := func(a *epath.ApplicationPath, i *AssemblyInstance) bool {
		if a == nil && i == nil {
			return true
		}
		return a.GetInstanceIDOrConnectionPoint() == i.InstanceID
	}

	for i, cp := range self.ConnectionPoints {
		if cp.Transport != freq.TransportClassAndTrigger {
			continue
		}
		if !doesMatch(config, cp.Config) {
			continue
		}
		if !doesMatch(prod, cp.Input) {
			continue
		}
		if !doesMatch(cons, cp.Output) {
			continue
		}
		return &self.ConnectionPoints[i], nil
	}

	return nil, fmt.Errorf("no matching connection point found")
}

// ConfigSize, in bytes.
func (self *ConnectionPoint) ConfigSize(
	_ cm.TransportClassAndTrigger,
) cip.UINT {
	if self.Config != nil {
		return self.Config.GetSize()
	}
	return 0
}

// InputSize, in bytes with header (via. RealTimeFormat)
func (self *ConnectionPoint) InputSize(
	class cm.TransportClassAndTrigger,
) cip.UINT {
	if self.Input != nil {
		return self.Input.GetSize() + self.TtoOFormat.Size(class)
	}
	return 0
}

// OutputSize, in bytes with header (via. RealTimeFormat)
func (self *ConnectionPoint) OutputSize(
	class cm.TransportClassAndTrigger,
) cip.UINT {
	if self.Output != nil {
		return self.Output.GetSize() + self.OtoTFormat.Size(class)
	}
	return 0
}

// TriggerAndTransport Mask for the EDS file.
func (self *ConnectionPoint) TriggerAndTransportMaskString() string {
	m := self.TriggerAndTransporrtMask()
	return fmt.Sprintf("0x%08X", m)
}

// TriggerAndTransport Mask for the EDS file.
func (self *ConnectionPoint) TriggerAndTransporrtMask() cip.DWORD {
	v := cip.DWORD(0)
	tc := self.Transport.TransportClass()
	if tc == cm.Class0 {
		v |= 1 << 0
	}
	if tc == cm.Class1 {
		v |= 1 << 1
	}
	if tc == cm.Class2 {
		v |= 1 << 2
	}
	if tc == cm.Class3 {
		v |= 1 << 3
	}
	pt := self.Transport.ProductionTrigger()
	if pt == cm.Cyclic {
		v |= 1 << 16
	}
	if pt == cm.ChangeOfState {
		v |= 1 << 17
	}
	if pt == cm.ApplicationObject {
		v |= 1 << 18
	}

	if self.Type&ListenOnly != 0 {
		v |= 1 << 24
	}
	if self.Type&InputOnly != 0 {
		v |= 1 << 25
	}
	if self.Type&ExclusiveOwner != 0 {
		v |= 1 << 26
	}
	if self.Type&RedundantOwner != 0 {
		v |= 1 << 27
	}
	return v
}

func (self *ConnectionPoint) ConnectionParametersString() string {
	return fmt.Sprintf("0x%08X", self.ConnectionParameters())
}

// ConnectionParameters for the EDS file.
func (self *ConnectionPoint) ConnectionParameters() cip.DWORD {
	v := cip.DWORD(0)

	v |= 1 << 0 // OtoT Fixed Size
	v |= 1 << 2 // TtoO Fixed Size

	switch self.OtoTFormat {
	case ModelessFormat:
		v |= 0 << 8
	case ZeroLengthFormat:
		v |= 1 << 8
	case HeartbeatFormat:
		v |= 3 << 8
	case Header32BitFormat:
		v |= 4 << 8
	}

	switch self.TtoOFormat {
	case ModelessFormat:
		v |= 0 << 12
	case ZeroLengthFormat:
		v |= 1 << 12
	case HeartbeatFormat:
		v |= 3 << 12
	case Header32BitFormat:
		v |= 4 << 12
	}

	v |= 1 << 17 // OtoT multicast
	v |= 1 << 18 // OtoT point2point
	v |= 1 << 21 // TtoO multicast
	v |= 1 << 22 // TtoO point2point

	v |= 1 << 26 // OtoT Scheduled
	v |= 1 << 30 // TtoO Scheduled.

	return v
}

// ConnectionPathString for EDS file
func (self *ConnectionPoint) ConnectionPathString() string {
	// Make EPATH string.
	// This is probably a special case, cleanup later.
	// The inverse is implemented elsewhere.
	data := []byte{
		0x20,
		cip.USINT(cip.AssemblyClassCode),
	}
	if self.Config != nil {
		data = append(data, 0x24)
		data = append(data, cip.USINT(self.Config.InstanceID))
	}
	if self.Output != nil {
		data = append(data, 0x2C)
		data = append(data, cip.USINT(self.Output.InstanceID))
	}
	if self.Input != nil {
		data = append(data, 0x2C)
		data = append(data, cip.USINT(self.Input.InstanceID))
	}

	s := ""
	for _, b := range data {
		if len(s) > 0 {
			s += " "
		}
		s += fmt.Sprintf("%02X", b)
	}
	return s
}
