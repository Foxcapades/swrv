package swrv

import "net/http"

// A Response builds the HTTP response that will be sent to the HTTP client that
// made the source request.
type Response interface {

	// GetCode returns the HTTP status code set on this Response instance.
	//
	// If unset, the Response status code defaults to 200.
	GetCode() int

	// WithCode sets the HTTP status code of the response to the given value.
	//
	// If unset, the Response status code defaults to 200.
	WithCode(code int) Response

	// GetBody returns the body for this Response.
	GetBody() any

	// WithBody sets the body on this Response instance.
	//
	// The body may be any object and, if not an io.Reader, will be passed to the
	// first matching ObjectSerializer registered with the Server.
	//
	// If the body is an io.ReadCloser, the body will be closed automatically by
	// the server after the response has been written to the client.
	WithBody(body any) Response

	// GetHeaders returns the ResponseHeaders attached to this Response instance.
	GetHeaders() ResponseHeaders

	// WithHeader sets the target response header to the given values.
	WithHeader(header, value string, values ...string) Response
}

// NewResponse creates a new Response instance.
//
// The created response will have no body or headers, and will default to a
// status code of 200.
func NewResponse() Response {
	return &response{
		code:    200,
		body:    nil,
		headers: responseHeaders{make(http.Header)},
	}
}

type response struct {
	code    int
	body    any
	headers responseHeaders
}

func (r *response) GetCode() int {
	return r.code
}

func (r *response) WithCode(code int) Response {
	r.code = code
	return r
}

func (r *response) GetBody() any {
	return r.body
}

func (r *response) WithBody(body any) Response {
	r.body = body
	return r
}

func (r *response) GetHeaders() ResponseHeaders {
	return r.headers
}

func (r *response) WithHeader(header, value string, values ...string) Response {
	r.headers.Set(header, value, values...)
	return r
}
