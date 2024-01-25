package cip

type AdditionalStatus = UINT

const (
	// See Volume 1: Table 3.5-37 "Connection Manager Service Error Codes"
	DuplicateForwardOpen                AdditionalStatus = 0x100
	OwnershipConflict                                    = 0x106
	TargetConnectionNotFound                             = 0x107
	InvalidConnectionParam                               = 0x108
	InvalidConnectionSize                                = 0x109
	TargetNotConfigured                                  = 0x110
	RPINotSupported                                      = 0x111
	RPINotAccepted                                       = 0x112
	OutOfConnections                                     = 0x113
	VendorIDMismatch                                     = 0x114 // or product type.
	DeviceTypeMismatch                                   = 0x115
	RevisionMismatch                                     = 0x116
	InvalidApplicationPath                               = 0x118
	NonListenOnlyConnectionNotOpened                     = 0x119
	TargetOutOfConnections                               = 0x11A
	TransportClassNotSupported                           = 0x11C
	TtoOTriggerNotSupported                              = 0x11D
	DirectionNotSupported                                = 0x11E
	InvalidOtoTConnectionType                            = 0x123
	InvalidTtoOConnectionType                            = 0x124
	InvalidConfigSize                                    = 0x126
	InvalidOtoTNetworkConnectionSize                     = 0x127
	InvalidTtoONetworkConnectionSize                     = 0x128
	InvalidConfigurationApplicationPath                  = 0x129
	InvalidConsumingApplicationPath                      = 0x12A
	InvalidProducingApplicationPath                      = 0x12B
	NullForwardOpenNotSupported                          = 0x132
	SerialNumberMismatch                                 = 0x13A
	InvalidConnectionPath                                = 0x315
	ParameterErrorUnconnectedService                     = 0x205 // for invalid sockaddr info.
)
