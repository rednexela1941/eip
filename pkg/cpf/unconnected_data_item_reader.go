package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/mr"
)

// Volume 2: 2-5.3.1
type UnconnectedDataItemReader interface {
	IReader
	//	NOTE: The format of the “data” field is dependent on the encapsulated protocol.  When used
	//
	// to encapsulate CIP, the format of the “data” field is that of a Message Router request or
	// Message Router reply.  See chapter 3 of this specification for details of the encapsulation of
	// UCMM messages.  See Volume 1, Chapter 2 for the format of the Message Router request and
	// reply packets.
	GetMessageRouterRequest() *mr.Request
}

type _UnconnectedDataItemReader struct {
	Item
	*mr.Request
}

func (self *_UnconnectedDataItemReader) GetMessageRouterRequest() *mr.Request {
	return self.Request
}

func (self *Item) ToUnconnectedDataItemReader(r bbuf.Reader) (UnconnectedDataItemReader, error) {
	item := &_UnconnectedDataItemReader{Item: *self}
	data, _ := item._getData(r)
	item.Data = data

	if r.Error() != nil {
		return nil, r.Error()
	}

	req, err := mr.NewRequest(item.Data)
	if err != nil {
		return item, err
	}
	item.Request = req

	return item, nil
	// return item, r.Error()
}
