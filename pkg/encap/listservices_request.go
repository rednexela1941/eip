package encap

/*
The ListServices command shall determine which encapsulation service classes the target
device supports. The ListServices command does not require that a session be established.
NOTE: Each service class has a unique type code, and an optional ASCII name.
*/

// Volume 2: Table 2-4.12
type ListServicesRequest Request

type _ListServicesRequest Packet

func (self *Packet) ToListServicesRequest() (ListServicesRequest, error) {
	return (*_ListServicesRequest)(self), nil
}
