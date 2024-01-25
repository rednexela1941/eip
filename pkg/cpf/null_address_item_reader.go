package cpf

import (
	"fmt"

	"github.com/rednexela1941/eip/pkg/bbuf"
)

// Volume 2: Table 2-5.4
type (
	NullAddressItemReader  IReader
	_NullAddressItemReader Item
)

func (self *Item) ToNullAddressItemReader(_ bbuf.Reader) (NullAddressItemReader, error) {
	if self.Length != 0 {
		return nil, fmt.Errorf("invalid length %d for %s", self.Length, self.TypeID.String())
	}
	return (*_NullAddressItemReader)(self), nil
}
