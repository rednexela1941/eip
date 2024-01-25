package cip

//go:generate stringer -type=ServiceCode
type ServiceCode USINT

// See Vol 1: Table A-3.1
const (
	GetAttributesAll             ServiceCode = 0x01
	SetAttributesAll             ServiceCode = 0x02
	GetAttributeList             ServiceCode = 0x03
	SetAttributeList             ServiceCode = 0x04
	Reset                        ServiceCode = 0x05
	Start                        ServiceCode = 0x06
	Stop                         ServiceCode = 0x07
	Create                       ServiceCode = 0x08
	Delete                       ServiceCode = 0x09
	MultipleServicePacket        ServiceCode = 0x0A
	ApplyAttributes              ServiceCode = 0x0D
	GetAttributeSingle           ServiceCode = 0x0E
	SetAttributeSingle           ServiceCode = 0x10
	FindNextObjectInstance       ServiceCode = 0x11
	Restore                      ServiceCode = 0x15
	Save                         ServiceCode = 0x16
	NoOP                         ServiceCode = 0x17
	GetMember                    ServiceCode = 0x18
	SetMember                    ServiceCode = 0x19
	InsertMember                 ServiceCode = 0x1A
	RemoveMember                 ServiceCode = 0x1B
	GroupSync                    ServiceCode = 0x1C
	GetConnectionPointMemberList ServiceCode = 0x1D

	// Table 3-5.10 Connection Manager Object-Specific Services
	ForwardClose         ServiceCode = 0x4E
	UnconnectedSend      ServiceCode = 0x52
	ForwardOpen          ServiceCode = 0x54
	GetConnectionData    ServiceCode = 0x56
	SearchConnectionData ServiceCode = 0x57
	GetConnectionOwner   ServiceCode = 0x5A
	LargeForwardOpen     ServiceCode = 0x5B

	GetAttributesAllResponse             ServiceCode = GetAttributesAll | 0x80
	SetAttributesAllResponse             ServiceCode = SetAttributesAll | 0x80
	GetAttributeListResponse             ServiceCode = GetAttributeList | 0x80
	SetAttributeListResponse             ServiceCode = SetAttributeList | 0x80
	ResetResponse                        ServiceCode = Reset | 0x80
	StartResponse                        ServiceCode = Start | 0x80
	StopResponse                         ServiceCode = Stop | 0x80
	CreateResponse                       ServiceCode = Create | 0x80
	DeleteResponse                       ServiceCode = Delete | 0x80
	MultipleServicePacketResponse        ServiceCode = MultipleServicePacket | 0x80
	ApplyAttributesResponse              ServiceCode = ApplyAttributes | 0x80
	GetAttributeSingleResponse           ServiceCode = GetAttributeSingle | 0x80
	SetAttributeSingleResponse           ServiceCode = SetAttributeSingle | 0x80
	FindNextObjectInstanceResponse       ServiceCode = FindNextObjectInstance | 0x80
	RestoreResponse                      ServiceCode = Restore | 0x80
	SaveResponse                         ServiceCode = Save | 0x80
	NoOPResponse                         ServiceCode = NoOP | 0x80
	GetMemberResponse                    ServiceCode = GetMember | 0x80
	SetMemberResponse                    ServiceCode = SetMember | 0x80
	InsertMemberResponse                 ServiceCode = InsertMember | 0x80
	RemoveMemberResponse                 ServiceCode = RemoveMember | 0x80
	GroupSyncResponse                    ServiceCode = GroupSync | 0x80
	GetConnectionPointMemberListResponse ServiceCode = GetConnectionPointMemberList | 0x80
	ForwardCloseResponse                 ServiceCode = ForwardClose | 0x80
	UnconnectedSendResponse              ServiceCode = UnconnectedSend | 0x80
	ForwardOpenResponse                  ServiceCode = ForwardOpen | 0x80
	GetConnectionDataResponse            ServiceCode = GetConnectionData | 0x80
	SearchConnectionDataResponse         ServiceCode = SearchConnectionData | 0x80
	GetConnectionOwnerResponse           ServiceCode = GetConnectionOwner | 0x80
	LargeForwardOpenResponse             ServiceCode = LargeForwardOpen | 0x80
)

func (self ServiceCode) IsGet() bool {
	return self == GetAttributeSingle || self == GetAttributesAll || self == GetAttributeList
}

func (self ServiceCode) IsSet() bool {
	return self == SetAttributeList || self == SetAttributeSingle || self == SetAttributesAll
}
