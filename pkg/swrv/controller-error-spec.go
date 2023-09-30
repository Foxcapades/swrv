package swrv

// ErrorControllerSpec defines a simplified ControllerSpec type which may be
// used to construct 404 or 405 error handlers.
type ErrorControllerSpec interface {
	// GetHandler returns the handler instance for the controller.
	GetHandler() RequestHandler

	// WithRequestFilters appends controller-specific request filters that will be
	// applied to incoming requests after the global request filters set on the
	// parent Server instance.
	WithRequestFilters(filters ...RequestFilter) ErrorControllerSpec

	// GetRequestFilters returns this controller's controller-specific
	// RequestFilter instances.
	GetRequestFilters() []RequestFilter

	// WithResponseFilters appends controller-specific response filters that will
	// be applied to outgoing responses after the global response filters set on
	// the parent Server instance.
	WithResponseFilters(filters ...ResponseFilter) ErrorControllerSpec

	// GetResponseFilters returns this controller's controller-specific
	// ResponseFilter instances.
	GetResponseFilters() []ResponseFilter
}

// NewErrorController constructs a new ErrorControllerSpec instance which may
// be used to construct an error handling controller for 404 or 405 errors.
func NewErrorController(handler RequestHandler) ErrorControllerSpec {
	return &errorControllerSpec{handler: handler}
}

type errorControllerSpec struct {
	in      []RequestFilter
	out     []ResponseFilter
	handler RequestHandler
}

func (c *errorControllerSpec) GetHandler() RequestHandler {
	return c.handler
}

func (c *errorControllerSpec) WithRequestFilters(filters ...RequestFilter) ErrorControllerSpec {
	c.in = append(c.in, filters...)
	return c
}

func (c *errorControllerSpec) GetRequestFilters() []RequestFilter {
	return c.in
}

func (c *errorControllerSpec) WithResponseFilters(filters ...ResponseFilter) ErrorControllerSpec {
	c.out = append(c.out, filters...)
	return c
}

func (c *errorControllerSpec) GetResponseFilters() []ResponseFilter {
	return c.out
}
