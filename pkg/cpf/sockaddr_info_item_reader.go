package cpf

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
)

// Volume 2: Table 2-5.9
type SockaddrInfoItemReader interface {
	IReader
	GetSockaddrInfo() SockaddrInfo
	IsTtoO() bool
	IsOtoT() bool
}

type _SockaddrInfoItemReader struct {
	Item
	SockaddrInfo
}

func (self *_SockaddrInfoItemReader) IsTtoO() bool                  { return self.TypeID == SockaddrInfoTtoO }
func (self *_SockaddrInfoItemReader) IsOtoT() bool                  { return self.TypeID == SockaddrInfoOtoT }
func (self *_SockaddrInfoItemReader) GetSockaddrInfo() SockaddrInfo { return self.SockaddrInfo }

func (self *Item) ToSockaddrInfoItemReader(r bbuf.Reader) (SockaddrInfoItemReader, error) {
	if self.Length != 16 {
		return nil, fmt.Errorf("invalid length %d for %s", self.Length, self.TypeID.String())
	}
	item := &_SockaddrInfoItemReader{Item: *self}
	r.Rb(&item.SockaddrInfo) // Big Endian.

	// fmt.Println(item.SockaddrInfo.SinAddr.ToIP())

	return item, r.Error()
}
