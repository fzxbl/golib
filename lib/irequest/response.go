package irequest

import (
	"io"
	"net/http"
)

type IResponse interface {
	GetStatus() string
	GetStatusCode() int
	GetBytesResp() []byte
	GetStringResp() string
	GetReaderResp() io.ReadSeeker
	GetHttpResp() http.Response
}

type response struct {
	status     string
	statusCode int
	rawContent []byte
	content    string
	httpResp   http.Response
	body       io.ReadSeeker
}

func (r response) GetStatus() string {
	return r.status
}
func (r response) GetStatusCode() int {
	return r.statusCode
}
func (r response) GetBytesResp() []byte {
	return r.rawContent
}
func (r response) GetStringResp() string {
	return r.content
}

func (r response) GetHttpResp() http.Response {
	return r.httpResp
}

func (r response) GetReaderResp() io.ReadSeeker {
	return r.body
}
