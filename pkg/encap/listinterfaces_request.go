package encap

/*
The optional List Interfaces command shall be used by a connection originator to identify non-
CIP communication interfaces associated with the target.  A session need not be established to
send this command.
*/

// Volume 2: Table 2-4.7
type (
	ListInterfacesRequest  Request
	_ListInterfacesRequest Packet
)

func (self *Packet) ToListInterfacesRequest() (ListInterfacesRequest, error) {
	return (*_ListInterfacesRequest)(self), nil
}
