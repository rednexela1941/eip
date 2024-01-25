package encap

/*
Either an originator or a target may send this command to terminate the session.  The receiver
shall initiate a close of the underlying TCP/IP connection when it receives this command.  The
session shall also be terminated when the transport connection between the originator and
target is terminated.  The receiver shall perform any other associated cleanup required on its
end.  There shall be no reply to this command, except in the event that the command is received via UDP.  If the command is received via UDP, the receiver shall reply with encapsulation
error code 0x01 (invalid or unsupported command).
*/

// Volume 2: Table 2-4.11
type UnregisterSessionRequest Request

type _UnregisterSessionRequest Packet

// The receiver shall not reject the UnRegisterSession due to unexpected values in the
// encapsulation header (invalid Session Handle, non-zero Status, non-zero Options, or additional command data).  In all cases the TCP connection shall be closed.
func (self *Packet) ToUnregisterSessionRequest() (UnregisterSessionRequest, error) {
	return (*_UnregisterSessionRequest)(self), nil
}
