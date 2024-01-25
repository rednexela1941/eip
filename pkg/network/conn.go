package network

type Conn interface {
	// net.Conn
	// TODO: fill out required fields, try to subclass net.Conn
	IsMulticast() bool
	IsUDP() bool
	IsTCP() bool
	LocalIP() IPv4
}
