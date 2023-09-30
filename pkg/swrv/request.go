package swrv

import (
	"io"
	"mime/multipart"
	"net/http"
)

// Request represents an incoming HTTP request.
type Request interface {

	// Raw returns the underling http.Request struct pointer.
	Raw() *http.Request

	// AdditionalContext returns the RequestContext object attached to this
	// request.
	//
	// Additional context may be used to store arbitrary data along with the
	// incoming request.
	AdditionalContext() RequestContext

	// GetHeader returns the first value associated with the given header name.
	//
	// If no such header was found on the request, the returned string will be
	// empty.
	//
	// This lookup is case-insensitive.
	GetHeader(header string) string

	// GetHeaders returns all the values associated with the given header name.
	//
	// If no such header was found on the request, the returned slice will be
	// empty.
	//
	// This lookup is case-insensitive.
	GetHeaders(header string) []string

	// HasQueryParam tests whether this request was sent with the target query
	// param as part of the URL.
	HasQueryParam(name string) bool

	// GetQueryParam fetches the value for the target query param from the request
	// URL.
	//
	// If the request URL did not contain the target query param, the returned
	// string will be empty.
	GetQueryParam(name string) string

	// GetCookie returns the cookie with the given name.
	//
	// If no such cookie was found, the return value will be nil.
	GetCookie(name string) *http.Cookie

	// HasCookie tests whether this request has the target cookie.
	HasCookie(name string) bool

	// GetCookies returns a slice of all the cookies sent with the request.
	GetCookies() []*http.Cookie

	// URIParam fetches the value of the URI param with the given name.
	//
	// URI params are the parameterized portion of a path set on the parent
	// controller.
	URIParam(name string) string

	// URIParams returns
	URIParams() map[string]string

	// Body returns an io.ReadCloser over the raw request body.
	Body() io.ReadCloser

	// WithBody executes the given function, passing in the request body.
	//
	// The request body will be automatically closed on return of the given
	// function.
	WithBody(fn func(reader io.Reader))

	// MultipartReader returns a MIME multipart reader if this is a
	// multipart/form-data or a multipart/mixed POST request, else returns nil and
	// an error.
	MultipartReader() (*multipart.Reader, error)
}
