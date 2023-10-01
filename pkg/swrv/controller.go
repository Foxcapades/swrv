package swrv

import (
	"bufio"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func newController(
	in []RequestFilter,
	out []ResponseFilter,
	hand RequestHandler,
	serial []ObjectSerializer,
	logger *logrus.Entry,
) http.Handler {
	return controller{
		inFilters:   in,
		outFilters:  out,
		handler:     hand,
		serializers: serial,
		logger:      logger,
	}
}

type controller struct {
	inFilters   []RequestFilter
	outFilters  []ResponseFilter
	handler     RequestHandler
	serializers []ObjectSerializer
	logger      *logrus.Entry
}

func (c controller) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	c.logger.Traceln("accepted request")
	var response Response

	// Attempt to close the request body (if it has one) once we're done
	// processing the request.
	if r.Body != nil {
		defer func(body io.ReadCloser) {
			if err := body.Close(); err != nil {
				c.logger.Errorln("failed to close request body")
			}
		}(r.Body)
	}

	request := WrapRequest(r)

	for _, in := range c.inFilters {
		if response = in.FilterRequest(request); response != nil {
			c.handleResponse(writer, request, response)
			return
		}
	}

	c.logger.Traceln("processed input filters, moving to request handler")

	if response = c.handler.HandleRequest(request); response != nil {
		c.handleResponse(writer, request, response)
		return
	}

	c.logger.Errorln("handler did not return a response")

	c.handleResponse(writer, request, newEmptyResponseError("request handler did not return a response, returning 500 error"))
}

func (c controller) handleResponse(writer http.ResponseWriter, request Request, response Response) {
	c.logger.Debugln("handling response")

	for _, out := range c.outFilters {
		if response = out.FilterResponse(request, response); response == nil {
			c.logger.Errorln("response filter did not return a response object, returning 500 error")
			response = newEmptyResponseError("response filter did not return a response")
		}
	}

	// Flag indicating whether we've already set a Content-Type header.  This is
	// used later when determining whether an ObjectSerializer should set a
	// Content-Type
	setContentType := false

	// Apply any response headers.
	response.GetHeaders().ForEach(func(header string, values []string) {
		for _, val := range values {
			if header == HeaderContentType {
				setContentType = true
			}
			writer.Header().Add(header, val)
		}
	})

	c.logger.Debugln("processing response body")

	// Fetch the response body.
	body := response.GetBody()

	// If there is no response body, then stop here.
	if body == nil {
		writer.WriteHeader(response.GetCode())
		c.logger.Traceln("response was nil, returning empty body")
		return
	}

	if reader, ok := body.(io.ReadCloser); ok {
		c.logger.Traceln("response body is a readcloser")

		writer.WriteHeader(response.GetCode())

		defer func(reader io.ReadCloser) {
			if err := reader.Close(); err != nil {
				c.logger.Errorln("failed to close body ReadCloser with error: ", err.Error())
			}
		}(reader)

		if _, err := bufio.NewWriter(writer).ReadFrom(reader); err != nil {
			c.logger.Errorln("failed to copy body from reader to response writer: " + err.Error())
		}

		return
	}

	if reader, ok := body.(io.Reader); ok {
		c.logger.Traceln("response body is a reader")

		writer.WriteHeader(response.GetCode())

		if _, err := bufio.NewWriter(writer).ReadFrom(reader); err != nil {
			c.logger.Errorln("failed to copy body from reader to response writer: " + err.Error())
		}

		return
	}

	// Lookup the matching object serializer
	var serializer ObjectSerializer
	for _, serial := range c.serializers {
		if serial.Matches(body) {
			serializer = serial
			break
		}
	}

	// If no matching object serializer was found, fallback to the default one.
	if serializer == nil {
		serializer = defaultObjectSerializer{}
	}

	// If the response didn't directly set a Content-Type header, set one now.
	if !setContentType {
		writer.Header().Set(HeaderContentType, serializer.ContentType())
	}

	writer.WriteHeader(response.GetCode())

	// Attempt to serialize the response body.
	serialized, err := serializer.Serialize(body)

	// If we failed to serialize the response body, fallback to a bad error.
	// TODO: handle this more gracefully?
	if err != nil {
		c.logger.Errorln("response body serialization failed with error: " + err.Error())
		writer.WriteHeader(500)
		writer.Header().Set(HeaderContentType, ContentTypeTextPlain)
		serialized = strings.NewReader("response body serialization failed!")
	}

	// If the serializer returned something closeable, then read it and close it.
	if reader, ok := serialized.(io.ReadCloser); ok {
		defer func(reader io.ReadCloser) {
			if err := reader.Close(); err != nil {
				c.logger.Errorln("failed to close body ReadCloser with error: ", err.Error())
			}
		}(reader)
	}

	if _, err := bufio.NewWriter(writer).ReadFrom(serialized); err != nil {
		c.logger.Errorln("failed to copy body from reader to response writer: " + err.Error())
	}
}
