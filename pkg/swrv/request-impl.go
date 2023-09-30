package swrv

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
)

// WrapRequest wraps the given http.Request pointer in a new Request instance.
//
// The new Request will have an empty RequestContext attached.
func WrapRequest(r *http.Request) Request {
	return &request{
		request: r,
		context: make(requestContext, 2),
	}
}

type request struct {
	request *http.Request
	context requestContext
}

func (r *request) Raw() *http.Request {
	return r.request
}

func (r *request) AdditionalContext() RequestContext {
	return r.context
}

// Body ////////////////////////////////////////////////////////////////////////

func (r *request) Body() io.ReadCloser {
	return r.request.Body
}

func (r *request) WithBody(fn func(io.Reader)) {
	defer r.request.Body.Close()
	fn(r.request.Body)
}

func (r *request) MultipartReader() (*multipart.Reader, error) {
	return r.request.MultipartReader()
}

// Query Params ////////////////////////////////////////////////////////////////

func (r *request) HasQueryParam(name string) bool {
	return r.request.URL.Query().Has(name)
}

func (r *request) GetQueryParam(name string) string {
	return r.request.URL.Query().Get(name)
}

// Headers /////////////////////////////////////////////////////////////////////

func (r *request) GetHeader(header string) string {
	return r.request.Header.Get(header)
}

func (r *request) GetHeaders(header string) []string {
	return r.request.Header.Values(header)
}

// Cookies /////////////////////////////////////////////////////////////////////

func (r *request) GetCookie(name string) *http.Cookie {
	if c, e := r.request.Cookie(name); e != nil {
		return nil
	} else {
		return c
	}
}

func (r *request) HasCookie(name string) bool {
	if _, e := r.request.Cookie(name); e != nil {
		return false
	} else {
		return true
	}
}

func (r *request) GetCookies() []*http.Cookie {
	return r.request.Cookies()
}

// URI Params //////////////////////////////////////////////////////////////////

func (r *request) URIParam(name string) string {
	return mux.Vars(r.request)[name]
}

func (r *request) URIParams() map[string]string {
	return mux.Vars(r.request)
}
