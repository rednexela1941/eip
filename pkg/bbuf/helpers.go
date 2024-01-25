package bbuf

import (
	"encoding/binary"
	"io"

	"github.com/rednexela1941/eip/pkg/cip"
)

func WShortString(w io.Writer, target string) error {
	s := target
	if len(s) > 0xff {
		s = s[:0xff]
	}
	l := cip.USINT(len(s))
	if err := Wl(w, l); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}

// Wl writes target in little endian to w.
func Wl(w io.Writer, target interface{}) error {
	return binary.Write(w, binary.LittleEndian, target)
}

// Wb writes target in big endian to w.
func Wb(w io.Writer, target interface{}) error {
	return binary.Write(w, binary.BigEndian, target)
}
