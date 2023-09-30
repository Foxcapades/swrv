package swrv

// A ResponseFilter is a filter that is applied to outgoing HTTP responses that
// may modify or replace the Response.
//
// ResponseFilters are applied in the order that they are registered, with
// global ResponseFilters, i.e. the ones set on the Server instance, being
// applied after the controller-specific ResponseFilters.
//
// All ResponseFilter instances will be called regardless of what they return.
type ResponseFilter interface {

	// FilterResponse may apply changes to, or replace entirely, the passed
	// Response instance.
	//
	// FilterResponse is expected to always return a response instance.  If it
	// returns nil, the Server will respond with a 500 error.
	FilterResponse(request Request, response Response) Response
}
