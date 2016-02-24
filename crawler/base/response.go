/* crawler scheduler response */
package base

import (
	"net/http"
)

type Response struct {
	httpResp      *http.Response
	primaryDomain string
	update        bool
	depth         uint32
}

func NewResponse(httpResp *http.Response, primaryDomain string, depth uint32) *Response {
	return &Response{
		httpResp:      httpResp,
		primaryDomain: primaryDomain,
		depth:         depth,
		update:        false,
	}
}

func NewUpdateResponse(httpResp *http.Response, primaryDomain string, depth uint32) *Response {
	return &Response{
		httpResp:      httpResp,
		primaryDomain: primaryDomain,
		depth:         depth,
		update:        true,
	}
}

func (resp *Response) HttpResp() *http.Response {
	return resp.httpResp
}

func (resp *Response) Depth() uint32 {
	return resp.depth
}

func (resp *Response) PrimaryDomain() string {
	return resp.primaryDomain
}

func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

func (resp *Response) Update() bool {
	return resp.update
}
