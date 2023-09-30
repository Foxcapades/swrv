package swrv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewServer(host string, port uint16) Server {
	return &server{
		logger: logrus.WithField("log-from", "server"),
		extras: &serverExtras{
			readTimeout:  30 * time.Second,
			writeTimeout: 30 * time.Second,
			host:         host,
			port:         port,
		},
	}
}

// A Server serves HTTP requests.
type Server interface {
	// WithLogger configures the server to use the given logrus logger instance
	// for internal server logging.
	//
	// This will overwrite any value previously set using WithLogger or
	// WithLoggerEntry.
	//
	// Passing nil will cause the server to panic with a nil pointer exception on
	// startup.
	WithLogger(logger *logrus.Logger) Server

	// WithLoggerEntry configures the server to use the given logrus logger entry
	// instance for internal server logging.
	//
	// This will overwrite any value previously set using WithLogger or
	// WithLoggerEntry.
	//
	// Passing nil will cause the server to panic with a nil pointer exception on
	// startup.
	WithLoggerEntry(logger *logrus.Entry) Server

	// WithReadTimeout sets the server's read timeout value to the given duration.
	//
	// If unset, the Server will default to a 30-second timeout.
	WithReadTimeout(timeout time.Duration) Server

	// WithWriteTimeout sets the server's write timeout value to the given
	// duration.
	//
	// If unset, the Server will default to a 30-second timeout.
	WithWriteTimeout(timeout time.Duration) Server

	// WithControllers adds the given ControllerSpec to the Server.
	//
	// When the Server is started, this ControllerSpec will be built into a
	// controller instance that will handle incoming HTTP requests that match the
	// target path and filters.
	WithControllers(controller ControllerSpec) Server

	// WithRequestFilters appends global RequestFilter instances that will be hit
	// for requests to any controller registered with the Server instance.
	//
	// Global RequestFilter instances are applied before controller specific
	// RequestFilter instances.
	WithRequestFilters(filters ...RequestFilter) Server

	// WithResponseFilters appends global ResponseFilter instances that will be
	// hit for outgoing responses from any controller registered with the Server
	// instance.
	//
	// Global ResponseFilter instances are applied after controller specific
	// ResponseFilter instances.
	WithResponseFilters(filters ...ResponseFilter) Server

	// WithObjectSerializers appends ObjectSerializer instances to the Server.
	//
	// ObjectSerializers are used to serialize non-stream objects into values that
	// may be streamed out to the requesting client.
	//
	// If no ObjectSerializer instances are provided, the default ObjectSerializer
	// will be used.  The default ObjectSerializer simply stringifies the object.
	//
	// ObjectSerializers are not applied to Response bodies of type io.Reader.
	//
	// ObjectSerializers will be tested in the order they are appended to the
	// server.  The first matching serializer will be used to serialize a Response
	// body.
	WithObjectSerializers(serializers ...ObjectSerializer) Server

	// With404Controller configures the Server's 404 Not Found controller, that
	// is, the controller that will be called when a client makes a request to an
	// endpoint that is not registered to the Server.
	//
	// Optionally requests to this controller may choose to use the global
	// RequestFilter and ResponseFilter instances like a normal controller.
	With404Controller(useGlobalFilters bool, controller ErrorControllerSpec) Server

	// With405Controller configures the Server's 405 Method Not Allowed
	// controller, that is, the controller that will be called with a client makes
	// a request to an endpoint using an HTTP method that is not supported by that
	// endpoint.
	//
	// Optionally requests to this controller may choose to use the global
	// RequestFilter and ResponseFilter instances like a normal controller.
	With405Controller(useGlobalFilters bool, controller ErrorControllerSpec) Server

	// Start starts the server, binding to the configured port and address,
	// optionally using a given router.
	//
	// If the router parameter is nil, a new router will be initialized for the
	// server.
	//
	// Once a server has started, no new filters, controllers, or serializers may
	// be registered.
	//
	// A server may only be started once.
	//
	// Example: No Router
	//   server := xhttp.NewServer(address, port)
	//   ...
	//   server.Start(nil)
	//
	// Example: With Router
	//
	//   router := mux.NewRouter()
	//   server := xhttp.NewServer(address, port)
	//   ...
	//   server.Start(router)
	Start(router *mux.Router)
}

type serverExtras struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	host         string
	port         uint16
	useFilt404   bool
	useFilt405   bool
}

type server struct {
	logger      *logrus.Entry
	started     bool
	inFilters   []RequestFilter
	outFilters  []ResponseFilter
	controllers []ControllerSpec
	serializers []ObjectSerializer
	handler404  ErrorControllerSpec
	handler405  ErrorControllerSpec
	extras      *serverExtras
}

// Logging /////////////////////////////////////////////////////////////////////

func (s *server) WithLogger(logger *logrus.Logger) Server {
	s.logger = logger.WithField("log-from", "server")
	return s
}

func (s *server) WithLoggerEntry(logger *logrus.Entry) Server {
	s.logger = logger.WithField("log-from", "server")
	return s
}

// Controller //////////////////////////////////////////////////////////////////

func (s *server) WithControllers(controller ControllerSpec) Server {
	if s.started {
		s.logger.Fatalln("cannot add controllers to a server after it has started")
	}
	s.controllers = append(s.controllers, controller)
	return s
}

// Filtering ///////////////////////////////////////////////////////////////////

func (s *server) WithRequestFilters(filters ...RequestFilter) Server {
	if s.started {
		s.logger.Fatalln("cannot add request filters to a server after it has started")
	}
	s.inFilters = append(s.inFilters, filters...)
	return s
}

func (s *server) WithResponseFilters(filters ...ResponseFilter) Server {
	if s.started {
		s.logger.Fatalln("cannot add response filters to a server after it has started")
	}
	s.outFilters = append(s.outFilters, filters...)
	return s
}

// Serialization ///////////////////////////////////////////////////////////////

func (s *server) WithObjectSerializers(serializers ...ObjectSerializer) Server {
	if s.started {
		s.logger.Fatalln("cannot set an object serializer on a server after it has started")
	}
	s.serializers = append(s.serializers, serializers...)
	return s
}

// Timeouts ////////////////////////////////////////////////////////////////////

func (s *server) WithReadTimeout(timeout time.Duration) Server {
	s.extras.readTimeout = timeout
	return s
}

func (s *server) WithWriteTimeout(timeout time.Duration) Server {
	s.extras.writeTimeout = timeout
	return s
}

// Error Handling //////////////////////////////////////////////////////////////

func (s *server) With404Controller(
	useGlobalFilters bool,
	controller ErrorControllerSpec,
) Server {
	s.handler404 = controller
	s.extras.useFilt404 = useGlobalFilters
	return s
}

func (s *server) With405Controller(
	useGlobalFilters bool,
	controller ErrorControllerSpec,
) Server {
	s.handler405 = controller
	s.extras.useFilt405 = useGlobalFilters
	return s
}

// Run /////////////////////////////////////////////////////////////////////////

func (s *server) Start(router *mux.Router) {
	if s.started {
		s.logger.Warnln("attempted to start a server instance more than once, ignoring")
	}

	if router == nil {
		s.logger.Debugln("no router passed to Start, using default router")
		router = mux.NewRouter()
	}

	s.build(router)

	if s.handler404 != nil {
		s.logger.Debugln("registering custom 404 handler")
		s.buildErrorController(s.extras.useFilt404, s.handler404, &router.NotFoundHandler, 404)
	}

	if s.handler405 != nil {
		s.logger.Debugln("registering custom 405 handler")
		s.buildErrorController(s.extras.useFilt405, s.handler405, &router.MethodNotAllowedHandler, 405)
	}

	serve := http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.extras.host, s.extras.port),
		Handler:      router,
		ReadTimeout:  s.extras.readTimeout,
		WriteTimeout: s.extras.writeTimeout,
	}

	s.clear()

	s.logger.Debugf("starting server at %s\n", serve.Addr)
	s.logger.Fatalln(serve.ListenAndServe())
}

// Internals ///////////////////////////////////////////////////////////////////

func (s *server) build(router *mux.Router) {
	if len(s.controllers) == 0 {
		s.logger.Fatalln("attempted to start a server with no controllers registered")
	}

	s.logger.Debugln("building controllers")
	for _, controller := range s.controllers {
		s.buildController(controller, router)
	}
}

func (s *server) clear() {
	s.inFilters = nil
	s.outFilters = nil
	s.controllers = nil
	s.serializers = nil
	s.handler405 = nil
	s.handler404 = nil
	s.extras = nil
}

func (s *server) buildErrorController(
	appendGlobals bool,
	spec ErrorControllerSpec,
	slot *http.Handler,
	code int,
) {
	var inFilters []RequestFilter
	var outFilters []ResponseFilter

	if appendGlobals {
		inFilters = append(s.inFilters, spec.GetRequestFilters()...)
		outFilters = append(spec.GetResponseFilters(), s.outFilters...)
	} else {
		inFilters = spec.GetRequestFilters()
		outFilters = spec.GetResponseFilters()
	}

	*slot = newController(
		inFilters,
		outFilters,
		spec.GetHandler(),
		s.serializers,
		s.logger.WithField("controller", code),
	)
}

func (s *server) buildController(spec ControllerSpec, router *mux.Router) {
	inFilters := append(s.inFilters, spec.GetRequestFilters()...)
	outFilters := append(spec.GetResponseFilters(), s.outFilters...)

	// Ensure we have a valid path
	if len(spec.GetPath()) == 0 {
		s.logger.Fatalln("controller has an empty path")
	}

	s.logger.Tracef("building controller %s\n", spec.GetPath())

	route := router.Path(spec.GetPath())

	// If the controller should only fire for specific HTTP methods
	if len(spec.GetMethods()) > 0 {
		route.Methods(spec.GetMethods()...)
	}

	// If the controller requires specific headers to be set.
	if len(spec.GetRequiredHeaders()) > 0 {
		pairs := make([]string, 0, len(spec.GetRequiredHeaders())*2)

		for head, match := range spec.GetRequiredHeaders() {
			pairs = append(pairs, head, match)
		}

		route.Headers(pairs...)
	}

	// Build the controller.
	route.Handler(newController(
		inFilters,
		outFilters,
		spec.GetHandler(),
		s.serializers,
		s.logger.WithField("controller", spec.GetPath()),
	))

}
