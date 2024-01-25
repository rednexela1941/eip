package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/mr"
)

// Volume 2: 2-5.3.2
type ConnectedDataItemReader interface {
	IReader
	GetSequenceCount() cip.UINT
	GetMessageRouterRequest() *mr.Request

	// NOTE: The format of the “data” field is dependent on the encapsulated protocol.  When used
	// to encapsulate CIP, the format of the “data” field is that of connected packet.  See chapter 3 of
	// this specification for details of the encapsulation of connected packets.  See chapter 3 of the
	// CIP Specification (Volume 1) for the format of connected packets.
}

type _ConnectedDataItemReader struct {
	Item
	SequenceCount cip.UINT
	*mr.Request
}

func (self *_ConnectedDataItemReader) GetMessageRouterRequest() *mr.Request { return self.Request }

func (self *_ConnectedDataItemReader) GetSequenceCount() cip.UINT { return self.SequenceCount }

func (self *Item) ToConnectedDataItemReader(r bbuf.Reader) (ConnectedDataItemReader, error) {
	item := &_ConnectedDataItemReader{Item: *self}
	data, _ := item._getData(r)
	item.Data = data
	if r.Error() != nil {
		return nil, r.Error()
	}

	b := bbuf.New(item.Data)
	b.Rl(&item.SequenceCount) // also optional, could h
	if b.Error() != nil {
		return nil, b.Error()
	}

	req, err := mr.NewRequest(b.Bytes())
	if err != nil {
		// TODO: don't really have to error out here.
		// just leave item.Request = nil
		return item, err
	}

	item.Request = req

	return item, nil
}
