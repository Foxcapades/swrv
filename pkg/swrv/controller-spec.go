package swrv

// NewController returns a new ControllerSpec instance which may be used to
// construct a controller for handling HTTP requests.
func NewController(path string, handler RequestHandler) ControllerSpec {
	return &controllerSpec{path: path, handler: handler, headers: make(map[string]string)}
}

// ControllerSpec defines a specification for a controller that will be built
// for the parent Server.
type ControllerSpec interface {
	// GetPath returns the URL path for the controller.
	GetPath() string

	// GetHandler returns the handler instance for the controller.
	GetHandler() RequestHandler

	// WithRequestFilters appends controller-specific request filters that will be
	// applied to incoming requests after the global request filters set on the
	// parent Server instance.
	WithRequestFilters(filters ...RequestFilter) ControllerSpec

	// GetRequestFilters returns this controller's controller-specific
	// RequestFilter instances.
	GetRequestFilters() []RequestFilter

	// WithResponseFilters appends controller-specific response filters that will
	// be applied to outgoing responses after the global response filters set on
	// the parent Server instance.
	WithResponseFilters(filters ...ResponseFilter) ControllerSpec

	// GetResponseFilters returns this controller's controller-specific
	// ResponseFilter instances.
	GetResponseFilters() []ResponseFilter

	// ForMethods sets the HTTP methods that the built controller will listen for.
	//
	// If set, the controller will only be called for matching HTTP methods.
	//
	// If unset, the controller will be called for any HTTP method.
	ForMethods(methods ...string) ControllerSpec

	// GetMethods returns the list of HTTP methods that the controller will listen
	// for.
	GetMethods() []string

	// GetRequiredHeaders returns the header requirements for this controller.
	//
	// The built controller will only be called for requests containing headers
	// that match the given requirements.
	GetRequiredHeaders() map[string]string

	// WithRequiredHeader sets a header requirement for request matching.  Only
	// requests that contain the given header set to the given value will be
	// matched and forwarded to this controller.
	//
	// If the given value string is empty, the matcher will match any value set
	// on the target header.
	WithRequiredHeader(header, value string) ControllerSpec
}

type controllerSpec struct {
	path    string
	methods []string
	in      []RequestFilter
	out     []ResponseFilter
	handler RequestHandler
	headers map[string]string
}

func (c *controllerSpec) GetPath() string {
	return c.path
}

func (c *controllerSpec) GetHandler() RequestHandler {
	return c.handler
}

func (c *controllerSpec) WithRequestFilters(filters ...RequestFilter) ControllerSpec {
	c.in = append(c.in, filters...)
	return c
}

func (c *controllerSpec) GetRequestFilters() []RequestFilter {
	return c.in
}

func (c *controllerSpec) WithResponseFilters(filters ...ResponseFilter) ControllerSpec {
	c.out = append(c.out, filters...)
	return c
}

func (c *controllerSpec) GetResponseFilters() []ResponseFilter {
	return c.out
}

func (c *controllerSpec) ForMethods(methods ...string) ControllerSpec {
	c.methods = append(c.methods, methods...)
	return c
}

func (c *controllerSpec) GetMethods() []string {
	return c.methods
}

func (c *controllerSpec) WithRequiredHeader(header, value string) ControllerSpec {
	c.headers[header] = value
	return c
}

func (c *controllerSpec) GetRequiredHeaders() map[string]string {
	return c.headers
}
