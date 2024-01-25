package cip

//go:generate stringer -type=GeneralStatus
type GeneralStatus USINT

// General Status Codes: See Vol 1 Table B-1.1
const (
	StatusSuccess                          GeneralStatus = 0x00
	StatusCommunicationProblem             GeneralStatus = 0x01
	StatusResourceUnavailable              GeneralStatus = 0x02
	StatusInvalidParameterValue            GeneralStatus = 0x03
	StatusPathSegmentError                 GeneralStatus = 0x04
	StatusPathDestinationUnknown           GeneralStatus = 0x05
	StatusPartialTransfer                  GeneralStatus = 0x06
	StatusConnectionLost                   GeneralStatus = 0x07
	StatusServiceNotSupported              GeneralStatus = 0x08
	StatusInvalidAttributeValue            GeneralStatus = 0x09
	StatusAttributeListError               GeneralStatus = 0x0A
	StatusAlreadyInRequestedMode           GeneralStatus = 0x0B
	StatusObjectStateConflict              GeneralStatus = 0x0C
	StatusObjectAlreadyExists              GeneralStatus = 0x0D
	StatusAttributeNotSetable              GeneralStatus = 0x0E
	StatusPriveledgeViolation              GeneralStatus = 0x0F
	StatusDeviceStateConflict              GeneralStatus = 0x10
	StatusReplyDataTooLarge                GeneralStatus = 0x11
	StatusFragmentationOfPrimitiveValue    GeneralStatus = 0x12
	StatusNotEnoughData                    GeneralStatus = 0x13
	StatusAttributeNotSupported            GeneralStatus = 0x14
	StatusTooMuchData                      GeneralStatus = 0x15
	StatusObjectInstanceDoesNotExist       GeneralStatus = 0x16
	StatusFragmentationOutOfSequence       GeneralStatus = 0x17
	StatusNoStoredAttributeData            GeneralStatus = 0x18
	StatusStoreOperationFailure            GeneralStatus = 0x19
	StatusRoutingFailureReqPacketTooLarge  GeneralStatus = 0x1A
	StatusRoutingFailureResPacketTooLarge  GeneralStatus = 0x1B
	StatusMissingAttributeListData         GeneralStatus = 0x1C
	StatusInvalidAttributeList             GeneralStatus = 0x1D
	StatusEmbeddedServiceError             GeneralStatus = 0x1E
	StatusVendorSpecificError              GeneralStatus = 0x1F
	StatusInvalidParameter                 GeneralStatus = 0x20
	StatusWriteOnlyOnceValueAlreadyWritten GeneralStatus = 0x21
	StatusInvalidReply                     GeneralStatus = 0x22
	StatusBufferOverflow                   GeneralStatus = 0x23
	StatusMessageFormatError               GeneralStatus = 0x24
	StatusKeyFailureInPath                 GeneralStatus = 0x25
	StatusPathSizeInvalid                  GeneralStatus = 0x26
	StatusUnexpectedAttributeInList        GeneralStatus = 0x27
	StatusInvalidMemeberID                 GeneralStatus = 0x28
	StatusMemberNotSettable                GeneralStatus = 0x29
	StatusGroup2OnlyServerGeneralFailure   GeneralStatus = 0x2A
	StatusUknownModbusError                GeneralStatus = 0x2B
	StatusAttributeNotGettable             GeneralStatus = 0x2C
	StatusInstanceNotDeletable             GeneralStatus = 0x2D
)
