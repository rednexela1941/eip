// Package epath implements the relevant EPATH parsing as defined in Volume 1: Appendix C-1.
package epath

import (
	"fmt"
	"log"
	"strings"

	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

type (
	EPATH struct {
		IsPadded         bool // true if Padded EPATH.
		Data             []cip.BYTE
		ElectronicKey    *ElectronicKey
		ApplicationPaths []ApplicationPath
	}
	PaddedEPATH = EPATH
)

// Parse EPATH and populate ElectronicKey and ApplicationPath fields in struct.
// returns (number of words read, error)

func (self *EPATH) Parse() (nWords cip.USINT, err error) {
	totalLen := len(self.Data)
	if totalLen == 0 {
		return 0, fmt.Errorf("EPATH is empty")
	}

	firstSegment := Segment(self.Data[0])
	reader := bbuf.New(self.Data)

	getWordsRead := func() cip.USINT {
		bytesRead := totalLen - reader.Len()
		return cip.USINT(bytesRead / 2)
	}

	if firstSegment.IsElectronicKey() {
		reader.Rl(&firstSegment) // read off first byte.

		k, err := NewElectronicKey(reader)
		if err != nil {
			return getWordsRead(), err
		}
		self.ElectronicKey = k
	}

	paths, err := parseApplicationPaths(reader, self.IsPadded)
	if err != nil {
		return getWordsRead(), err
	}
	self.ApplicationPaths = paths
	return getWordsRead(), err
}

func New(data []cip.BYTE) *EPATH {
	return &EPATH{
		IsPadded: false,
		Data:     data,
	}
}

func NewPadded(data []cip.BYTE) *PaddedEPATH {
	return &PaddedEPATH{
		IsPadded: true,
		Data:     data,
	}
}

func (self *EPATH) HasDataSegment() bool {
	log.Println("HasDataSegment not implemeted")
	return false
}

func (self *EPATH) String() string {
	paths := make([]string, len(self.Data))
	for i, v := range self.Data {
		paths[i] = fmt.Sprintf("%02X", v)
	}
	return strings.Join(paths, " ")
}
