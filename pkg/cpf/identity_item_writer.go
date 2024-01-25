package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/identity"
)

// Volume 2: Table 2-4.4
type _CIPIdentityItem struct {
	ItemHeader
	ProtocolVersion cip.UINT
	SockaddrInfo
	identity.Identity
}

// AddIdentityItem adds a new CIP identity item to the CPF packet Writer.
func (self *_Writer) AddIdentityItem(sockaddrInfo *SockaddrInfo, id *identity.Identity, protocolVersion cip.UINT) {
	self.Items = append(self.Items, &_CIPIdentityItem{
		ItemHeader: ItemHeader{
			TypeID: CIPIdentity,
		},
		SockaddrInfo:    *sockaddrInfo,
		Identity:        *id,
		ProtocolVersion: protocolVersion,
	})
}

func (self *_CIPIdentityItem) WriteTo(w bbuf.Writer) error {
	temp := bbuf.New(nil)
	temp.Wl(self.ProtocolVersion)
	temp.Wb(&self.SockaddrInfo)

	if err := self.Identity.WriteTo(temp); err != nil {
		return err
	}
	if temp.Error() != nil {
		return temp.Error()
	}
	data := temp.Bytes()

	self.Length = cip.UINT(len(data))

	w.Wl(&self.ItemHeader)
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}
