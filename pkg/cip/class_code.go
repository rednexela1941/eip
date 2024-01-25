package cip

//go:generate stringer -type=ClassCode
type ClassCode UINT

const (
	IdentityClassCode            ClassCode = 0x01
	MessageRouterClassCode       ClassCode = 0x02
	DeviceNetClassCode           ClassCode = 0x03
	AssemblyClassCode            ClassCode = 0x04
	ConnectionClassCode          ClassCode = 0x05
	ConnectionManagerClassCode   ClassCode = 0x06
	RegisterClassCode            ClassCode = 0x07
	DiscreteInputPointClassCode  ClassCode = 0x08
	DiscreteOutputPointClassCode ClassCode = 0x09
	AnalogInputPointClassCode    ClassCode = 0x0A
	AnalogOutputPointClassCode   ClassCode = 0x0B

	// ...
	PortClassCode           ClassCode = 0xF4
	TCPIPInterfaceClassCode ClassCode = 0xF5
	EthernetLinkClassCode   ClassCode = 0xF6

	// ...
	LLDPManagementClassCode ClassCode = 0x109
	LLDPDataTableClassCode  ClassCode = 0x10A
)
