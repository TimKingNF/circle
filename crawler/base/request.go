/* crawler scheduler request */
package base

import (
	"net/http"
)

type Request struct {
	httpReq       *http.Request
	primaryDomain string

	// 代表该 Request 的当前深度值，一旦大于 crawler 指定的最大深度值 该请求将会被忽略
	depth uint32

	maxDepth uint32
}

func NewRequest(httpReq *http.Request, primaryDomain string, depth uint32) *Request {
	return &Request{
		httpReq:       httpReq,
		primaryDomain: primaryDomain,
		depth:         depth,
	}
}

func NewMaxRequest(httpReq *http.Request, primaryDomain string, depth, maxDepth uint32) *Request {
	return &Request{
		httpReq:       httpReq,
		primaryDomain: primaryDomain,
		depth:         depth,
		maxDepth:      maxDepth,
	}
}

func (req *Request) HttpReq() *http.Request {
	return req.httpReq
}

func (req *Request) Depth() uint32 {
	return req.depth
}

func (req *Request) MaxDepth() uint32 {
	return req.maxDepth
}

func (req *Request) PrimaryDomain() string {
	return req.primaryDomain
}

func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}
