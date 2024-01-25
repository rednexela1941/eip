package adapter

import (
	"log"

	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/encap"
	"github.com/rednexela1941/eip/pkg/mr"
)

type (
	Request struct {
		*mr.Request
		Parent encap.Request
		*RequestContext
	}

	Response interface {
		mr.ResponseWriter
		Parent() encap.Reply
	}
)

type _Response struct {
	mr.ResponseWriter
	parent encap.Reply
}

func NewRequest(parent encap.Request, ctx *RequestContext, req *mr.Request) *Request {
	return &Request{
		Request:        req,
		Parent:         parent,
		RequestContext: ctx,
	}
}

func NewResponse(parent encap.Reply, response mr.ResponseWriter) *_Response {
	return &_Response{
		parent:         parent,
		ResponseWriter: response,
	}
}

func (self *_Response) Parent() encap.Reply {
	return self.parent
}

func (self *_Adapter) Route(req *Request, res Response) {
	self._Route(req, res)

	if res.Error() != nil {
		self.Logger.Fatal(res.Error())
	}

	pathStr := "UnknownPath"
	if len(req.RequestPath.ApplicationPaths) > 0 {
		pathStr = req.RequestPath.ApplicationPaths[0].DebugString()
	}

	self.Logger.Printf(
		"\t%s %s %s [%d bytes]",
		req.Service.String(),
		pathStr,
		res.PeekGeneralStatus().String(),
		res.Len(),
	)
}

func (self *_Adapter) _Route(req *Request, res Response) {
	res.SetReplyService(req.Service | 0x80)

	n, err := req.RequestPath.Parse()
	if err != nil {
		log.Println(err, n)
		res.SetGeneralStatus(cip.StatusPathSegmentError)
		return
	}

	paths := req.RequestPath.ApplicationPaths
	if len(paths) == 0 {
		self.Logger.Println("no paths found", len(paths))
		res.SetGeneralStatus(cip.StatusPathSegmentError)
		return
	}

	if len(paths) > 1 {
		self.Logger.Println("too many paths", len(paths))
		res.SetGeneralStatus(cip.StatusPathSegmentError)
		return
	}

	if req.RequestPath.ElectronicKey != nil {
		// validate the electronic key.

		valid, addStatus := ValidateElectronicKey(
			self.Identity,
			req.RequestPath.ElectronicKey,
		)
		if !valid {
			res.SetGeneralStatus(cip.StatusKeyFailureInPath)
			for _, s := range addStatus {
				res.AddAdditionalStatus(s)
			}
			return
		}
	}

	path := paths[0]

	class, ok := self.Classes[path.ClassID]
	if !ok {
		res.SetGeneralStatus(cip.StatusObjectInstanceDoesNotExist)
		return
	}

	class.CallService(req, res)
}
