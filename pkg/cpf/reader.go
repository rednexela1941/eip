package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

type (
	// Volume 2: Table 2-5.1
	// Common Packet Format Reader interface.
	Reader interface {
		GetItemCount() cip.UINT
		GetItemWithType(ItemID) (IReader, bool)
		GetItem(index int) IReader
	}

	// Volume 2: Table 2-5.2
	// Common Packet Format Item Reader interface.
	IReader interface {
		GetTypeID() ItemID
		GetLength() cip.UINT
	}
)

type _Reader []IReader

func (self _Reader) GetItemCount() cip.UINT    { return cip.UINT(len(self)) }
func (self _Reader) GetItem(index int) IReader { return self[index] }
func (self _Reader) GetItemWithType(t ItemID) (IReader, bool) {
	for _, ir := range self {
		if ir.GetTypeID() == t {
			return ir, true
		}
	}
	return nil, false
}

func (self *Packet) GetItemCount() cip.UINT    { return self.ItemCount }
func (self *Packet) GetItem(index int) IReader { return &self.Items[index] }

func (self *Packet) GetItemWithType(t ItemID) (IReader, bool) {
	for i, it := range self.Items {
		if it.GetTypeID() == t {
			return &self.Items[i], true
		}
	}
	return nil, false
}

func (self *ItemHeader) GetTypeID() ItemID   { return self.TypeID }
func (self *ItemHeader) GetLength() cip.UINT { return self.Length }

func NewReader(buffer bbuf.Reader) (Reader, error) {
	return _NewReader(buffer, false)
}

func NewIOReader(buffer bbuf.Reader) (Reader, error) {
	return _NewReader(buffer, true)
}

func _NewReader(buffer bbuf.Reader, isIO bool) (Reader, error) {
	var numItems cip.UINT
	buffer.Rl(&numItems)

	r := make(_Reader, numItems)

	for i := 0; i < int(numItems); i++ {
		ir, err := NewIReader(buffer, isIO)
		if err != nil {
			return nil, err
		}
		r[i] = ir
	}

	return r, nil
}

func NewIReader(r bbuf.Reader, isIO bool) (IReader, error) {
	item := new(Item)
	r.Rl(&item.ItemHeader)
	if r.Error() != nil {
		return nil, r.Error()
	}
	switch item.TypeID {
	case AddressNull:
		return item.ToNullAddressItemReader(r)
	case AddressConnected:
		return item.ToConnectedAddressItemReader(r)
	case SequencedAddress:
		return item.ToSequencedAddressItemReader(r)
	case UnconnectedMessage:
		return item.ToUnconnectedDataItemReader(r)
	case ConnectedTransportPacket:
		if isIO {
			break // return generic item.
		}
		return item.ToConnectedDataItemReader(r)
	case SockaddrInfoOtoT, SockaddrInfoTtoO:
		return item.ToSockaddrInfoItemReader(r)
	default:
	}
	data, err := item._getData(r)
	if err != nil {
		return nil, err
	}
	item.Data = data
	return item, r.Error()
	// return nil, r.Error()

}
