// Package mr implements Message Router Request/Response types as specfied
// in Volume 1: 2-4
package mr

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/epath"
)

// See Volume 1: Table 2-4.1
type Request struct {
	Service              cip.ServiceCode    // Service code of the request
	RequestPathSizeWords cip.USINT          // The number of 16 bit words in the Request_Path field.
	RequestPath          *epath.PaddedEPATH // Shall be of the form [ElectronicKey], Application Path
	RequestData          []cip.OCTET
}

func NewRequest(data []byte) (*Request, error) {
	r := bbuf.New(data)
	req := new(Request)
	r.Rl(&req.Service)
	r.Rl(&req.RequestPathSizeWords)
	if r.Error() != nil {
		return nil, r.Error()
	}
	reqPathData := make([]cip.BYTE, 2*int(req.RequestPathSizeWords))
	r.Rl(&reqPathData)
	req.RequestPath = epath.NewPadded(reqPathData)
	req.RequestData = r.Bytes()
	return req, r.Error()
}
