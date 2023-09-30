package swrv

import "net/http"

// ResponseHeaders is a mapping of header values that will be sent out with an
// HTTP response.
type ResponseHeaders interface {
	// Set puts the given header value(s) into the header mapping, replacing any
	// value previously set for the target header.
	Set(header, value string, values ...string)

	// Append puts the given header value(s) into the header mapping, appending
	// to any value(s) previously set for the target header.
	Append(header, value string, values ...string)

	// GetFirst returns the first value set for the target header.
	GetFirst(header string) (value string, found bool)

	// GetAll returns all values set for the target header.
	GetAll(header string) (values []string, found bool)

	// GetNth returns the nth value set for the target header.
	GetNth(header string, n int) (value string, found bool)

	// ForEach iterates over all the set headers and calls the given function on
	// each.
	ForEach(fn func(header string, values []string))
}

type responseHeaders struct {
	head http.Header
}

func (r responseHeaders) Set(header string, value string, values ...string) {
	r.head.Set(header, value)
	for _, val := range values {
		r.head.Add(header, val)
	}
}

func (r responseHeaders) Append(header string, value string, values ...string) {
	r.head.Add(header, value)
	for _, val := range values {
		r.head.Add(header, val)
	}
}

func (r responseHeaders) GetFirst(header string) (string, bool) {
	if found, ok := r.head[header]; ok {
		if len(found) == 0 {
			return "", true
		} else {
			return found[0], true
		}
	} else {
		return "", false
	}
}

func (r responseHeaders) GetAll(header string) ([]string, bool) {
	if found, ok := r.head[header]; ok {
		return found, true
	} else {
		return nil, false
	}
}

func (r responseHeaders) GetNth(header string, n int) (string, bool) {
	if found, ok := r.head[header]; ok {
		if len(found) > n {
			return found[n], true
		}
	}

	return "", false
}

func (r responseHeaders) ForEach(fn func(header string, values []string)) {
	for k, v := range r.head {
		fn(k, v)
	}
}
