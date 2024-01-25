package encap

import "github.com/rednexela1941/eip/pkg/bbuf"

// Volume 2: Table 2-4.8
type ListInterfacesReply interface {
	Reply
	bbuf.Writer
}

func NewListInterfacesReply() ListInterfacesReply {
	return _NewReply()
}
