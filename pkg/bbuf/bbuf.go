// Package bbuf implements a bytes.Buffer derived utility class for reading and writing bytes with appropriate endianess.
package bbuf

import (
	"bytes"
	"encoding/binary"
)

type (
	Errorable interface {
		Error() error
	}

	Reader interface {
		Errorable
		Rl(interface{}) // read little-endian from buffer to interface.
		Rb(interface{}) // read big-endian from buffer to interface.
		Read([]byte) (int, error)
		Len() int
	}

	Writer interface {
		Errorable
		Wl(interface{}) // write little-endian to buffer from interface.
		Wb(interface{}) // write big-endian to buffer from interface.
		Write([]byte) (int, error)
		Len() int
	}

	ReadWriter interface {
		Reader
		Writer
	}
)

type (
	_Buffer struct {
		*bytes.Buffer
		err error
	}
	Buffer = _Buffer
)

func New(data []byte) *_Buffer {
	var buffer *bytes.Buffer
	if data == nil {
		buffer = new(bytes.Buffer)
	} else {
		buffer = bytes.NewBuffer(data)
	}
	return &_Buffer{
		Buffer: buffer,
	}
}

func (self *_Buffer) updateErr(err error) {
	if err != nil {
		self.err = err
	}
}

func (self *_Buffer) Error() error { return self.err }

func (self *_Buffer) Rl(target interface{}) {
	self.updateErr(binary.Read(self, binary.LittleEndian, target))
}

func (self *_Buffer) Rb(target interface{}) {
	self.updateErr(binary.Read(self, binary.BigEndian, target))
}

func (self *_Buffer) Wl(target interface{}) {
	self.updateErr(binary.Write(self, binary.LittleEndian, target))
}

func (self *_Buffer) Wb(target interface{}) {
	self.updateErr(binary.Write(self, binary.BigEndian, target))
}

// func (self *_Buffer) Write(data []byte) (int, error) {
// 	return self.Write(data)
// }
