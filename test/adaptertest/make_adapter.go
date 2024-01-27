package adaptertest

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/adapter"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/cm"
	"github.com/rednexela1941/eip/pkg/identity"
)

// silence warning.
func noop(a interface{}) { return }

const (
	GetSingle = adapter.GetSingle
	GetAll    = adapter.GetAll
)

type Response = adapter.Response

func testSetup() {
	// inputAssm := a.AddAssemblyInstance(100, 178*cip.BYTESize)
	//
	// inputAssm.AddUINTParam(
	// 	"Estop", &uint_ptr,
	// ).SetHelpString(
	// 	"Help String"
	// ).SetUnitString(
	// 	"Interface"
	// ).SetMinumumValue(
	// 	0
	// ).SetMaximumValue(
	// 	1
	// ).SetDefaultValue(
	// 	0
	// ).OnGet(
	// 	func (w bbuf.Writer) { w.Wl(true) }
	// ).OnSet(
	// 	func (w bbuf.Writer) { w.Rl(false) }
	// )
	// inputAssm.AddParam("Estop")
}

func makeDefaultAssembly(inst *adapter.AssemblyInstance, size int) {
	for i := 0; i < size; i++ {
		name := fmt.Sprintf("Inst %d Param %d", inst.InstanceID, i)

		inst.AddBOOLParam(name).SetHelpString(
			"help string here...",
		).OnGet(
			func() cip.BOOL { return true },
		).OnSet(func(v cip.BOOL) {
			fmt.Printf("%s set to %t\n", name, v)
		})

	}
}

func InitAssemblies(a *adapter.Adapter) {
	inputAssm := a.AddAssemblyInstance("Input Assembly", 100)
	makeDefaultAssembly(inputAssm, 178)

	outputAssm := a.AddAssemblyInstance("Output Assembly", 150)
	makeDefaultAssembly(outputAssm, 184)

	configAssm := a.AddAssemblyInstance("Config Assembly", 151)
	hbInput := a.AddAssemblyInstance("Heartbeat Input Assembly", 152)
	hbListen := a.AddAssemblyInstance("Heartbeat Listen Assembly", 153)

	// setupDummyFuncs(inputAssm)
	// setupDummyFuncs(outputAssm)
	// setupDummyFuncs(configAssm)
	// setupDummyFuncs(hbInput)
	// setupDummyFuncs(hbListen)

	a.AddConnectionPoint(
		"Exclusive Owner",
		cm.Class1|cm.Cyclic|cm.DirectionClient,
		adapter.ExclusiveOwner,
		configAssm,
		inputAssm,
		outputAssm,
		adapter.ModelessFormat,
		adapter.ModelessFormat,
	)
	a.AddConnectionPoint(
		"Input Only",
		cm.Class1|cm.Cyclic|cm.DirectionClient,
		adapter.InputOnly,
		configAssm,
		inputAssm,
		hbInput,
		adapter.HeartbeatFormat,
		adapter.ModelessFormat,
	)
	a.AddConnectionPoint(
		"Listen Only",
		cm.Class1|cm.Cyclic|cm.DirectionClient,
		adapter.ListenOnly,
		configAssm,
		inputAssm,
		hbListen,
		adapter.HeartbeatFormat,
		adapter.ModelessFormat,
	)

	noop(inputAssm)
	noop(outputAssm)
	noop(configAssm)
	noop(hbInput)
	noop(hbListen)
}

func MakeTestAdapter() *adapter.Adapter {
	i := &identity.Identity{
		VendorID:    1,
		DeviceType:  43,
		ProductCode: 65001,
		Revision: identity.Revision{
			Major: 1,
			Minor: 4,
		},
		SerialNumber: 0x183fb,
		ProductName:  "Linux Box",
	}
	a := adapter.New(i)
	InitAssemblies(a)
	return a
}
