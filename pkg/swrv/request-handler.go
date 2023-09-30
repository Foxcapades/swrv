package swrv

// RequestHandler represents the core part of a controller that actually
// processes a request.
type RequestHandler interface {

	// HandleRequest is called to process incoming requests and transform them
	// into Response objects to be returned the HTTP client caller.
	HandleRequest(request Request) Response
}

// A RequestHandlerFunc is a function that implements the RequestHandler
// interface.
type RequestHandlerFunc func(request Request) Response

func (r RequestHandlerFunc) HandleRequest(request Request) Response {
	return r(request)
}
