package mr

import (
	"github.com/rednexela1941/eip/pkg/bbuf"
	"github.com/rednexela1941/eip/pkg/cip"
)

type (
	ResponseHeader struct {
		ReplyService                cip.ServiceCode   // Service code of request + 0x80
		Reserved                    cip.OCTET         // shall equal 0
		GeneralStatus               cip.GeneralStatus // Volume 1: Appendix B
		SizeOfAdditionalStatusWords cip.USINT         // Number of 16 bit workds in the Additional Status array.
	}

	// See Volume 1: Table 2-4.2
	Response struct {
		ResponseHeader
		AdditionalStatus []cip.AdditionalStatus
		ResponseData     []cip.OCTET
	}

	ResponseWriter interface {
		PeekGeneralStatus() cip.GeneralStatus // return currently set general status code (for logging)

		SetReplyService(cip.ServiceCode)
		SetGeneralStatus(cip.GeneralStatus)
		AddAdditionalStatus(cip.AdditionalStatus)
		AddAdditionalStatusNoReplace(cip.AdditionalStatus)
		bbuf.Writer

		Encode() ([]byte, error)
	}
)

type _ResponseWriter struct {
	ResponseHeader
	AdditionalStatus []cip.AdditionalStatus
	*bbuf.Buffer     // ResponseData
}

// ToStruct converts a ResponseWriter to Response structure format.
func (self *_ResponseWriter) ToStruct() *Response {
	r := new(Response)
	r.ResponseHeader = self.ResponseHeader
	r.AdditionalStatus = self.AdditionalStatus
	r.ResponseData = self.Buffer.Bytes()
	return r
}

// AddAdditionalStatusToArray -- ignoring duplicate statuses.
func AddAdditionalStatusToArray(
	arr []cip.AdditionalStatus,
	status cip.AdditionalStatus,
) []cip.AdditionalStatus {
	for _, s := range arr {
		if s == status {
			return arr
		}
	}
	arr = append(arr, status)
	return arr
}

// TODO: implement this so that it doesn't replace additional status.
// this is for attaching the connection sizes
// in a failed forward open request (connection size mismatch)
// for the time being, it will stay the same.
func (self *_ResponseWriter) AddAdditionalStatusNoReplace(addStatus cip.AdditionalStatus) {
	self.AddAdditionalStatus(addStatus)
}

func (self *_ResponseWriter) AddAdditionalStatus(addStatus cip.AdditionalStatus) {
	// will not add duplicates

	self.AdditionalStatus = AddAdditionalStatusToArray(self.AdditionalStatus, addStatus)
}

func (self *ResponseHeader) PeekGeneralStatus() cip.GeneralStatus { return self.GeneralStatus }

func (self *ResponseHeader) SetGeneralStatus(status cip.GeneralStatus) {
	self.GeneralStatus = status
}

func (self *ResponseHeader) SetReplyService(service cip.ServiceCode) {
	self.ReplyService = service
}

func NewResponseWriter() *_ResponseWriter {
	r := new(_ResponseWriter)
	r.Buffer = bbuf.New(nil)
	return r
}

func (self *_ResponseWriter) Encode() ([]byte, error) {
	self.SizeOfAdditionalStatusWords = cip.USINT(len(self.AdditionalStatus))
	buffer := bbuf.New(nil)
	buffer.Wl(self.ResponseHeader)
	buffer.Wl(self.AdditionalStatus)
	if _, err := buffer.Write(self.Buffer.Bytes()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), buffer.Error()
}

func NewResponseFrom(data []byte) (*Response, error) {
	r := bbuf.New(data)
	res := new(Response)
	r.Rl(&res.ResponseHeader)
	if r.Error() != nil {
		return nil, r.Error()
	}
	res.AdditionalStatus = make([]cip.WORD, res.SizeOfAdditionalStatusWords)
	r.Rl(&res.AdditionalStatus)
	res.ResponseData = r.Bytes()
	return res, r.Error()
}
