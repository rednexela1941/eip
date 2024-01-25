package cpf

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/identity"
)

type (
	// Volume 2: Table 2-5.1
	// Common Packet Format Writer interface.
	Writer interface {
		AddItem() IWriter

		AddNullAddressItem()
		AddConnectedAddressItem(connID cip.UDINT)
		AddIdentityItem(*SockaddrInfo, *identity.Identity, cip.UINT /* protocol version */)
		AddConnectedDataItem(sequenceCount cip.UINT) ConnectedDataItemWriter
		AddUnconnectedDataItem() UnconnectedDataItemWriter
		AddSockaddrInfoItem(ItemID, SockaddrInfo)
		AddSequencedAddressItem(connID cip.UDINT, seqNo cip.UDINT)
		AddIODataItem() IWriter

		Encode() ([]byte, error)
	}

	// Volume 2: Table2-5.2
	// Common Packet Format Item Writer interface.
	IWriter interface {
		IHeaderWriter
		bbuf.Writer
	}

	IHeaderWriter interface {
		SetTypeID(ItemID)

		WriteTo(bbuf.Writer) error
	}
)

type (
	// internal packet writer struct.
	_Writer struct {
		Items []IHeaderWriter
	}

	// internal item writer struct
	_IWriter struct {
		ItemHeader
		*bbuf.Buffer
	}
)

// AddItem adds a generic CPF Item to the CPF Writer.
func (self *_Writer) AddItem() IWriter {
	iw := new(_IWriter)
	iw.Buffer = bbuf.New(nil)
	self.Items = append(self.Items, iw)
	return iw
}

func (self *_IWriter) WriteTo(w bbuf.Writer) error {
	data := self.Bytes()
	self.Length = cip.UINT(len(data))
	w.Wl(self.ItemHeader)
	if len(data) > 0 {
		w.Wl(data)
	}
	return w.Error()
}

func (self *ItemHeader) SetTypeID(typeID ItemID) { self.TypeID = typeID }

func NewWriter() Writer {
	w := new(_Writer)
	w.Items = make([]IHeaderWriter, 0, 2)
	return w
}

func (self *_Writer) getItemCount() cip.UINT {
	return cip.UINT(len(self.Items))
}

func (self *_Writer) Encode() ([]byte, error) {
	buffer := bbuf.New(nil)
	buffer.Wl(self.getItemCount())
	for _, item := range self.Items {
		if err := item.WriteTo(buffer); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), buffer.Error()
}
