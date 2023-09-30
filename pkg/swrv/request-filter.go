package swrv

// A RequestFilter is a filter that is applied to incoming HTTP requests that
// may modify the RequestContext, or halt processing of the request by returning
// a non-nil Response.
//
// RequestFilters are applied in the order that they are registered with global
// RequestFilters, i.e. ones set on the Server instance, applied before
// controller-specific RequestFilters.
type RequestFilter interface {

	// FilterRequest may apply changes to the incoming Request's RequestContext,
	// or optionally, halt processing of the Request by returning a non-nil
	// Response.
	//
	// If this method returns a non-nil response, no further RequestFilters will
	// be called, and the RequestHandler will not be reached.  Instead, the Server
	// will return the response provided by this method.
	//
	// If this method returns a nil value, the request will be passed to the next
	// RequestFilter instance registered to either the server or the controller.
	FilterRequest(request Request) Response
}

// RequestFilterFunc defines a function that implements the RequestFilter
// interface.
type RequestFilterFunc func(request Request) Response

func (r RequestFilterFunc) FilterRequest(request Request) Response {
	return r(request)
}
