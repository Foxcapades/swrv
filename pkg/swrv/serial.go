package swrv

import "io"

// MatcherFn defines a function that may be used as a body object matcher in an
// ObjectSerializer.
type MatcherFn = func(object any) bool

// An ObjectSerializer is used to serialize response objects into a stream that
// may be passed to the HTTP client caller.
type ObjectSerializer interface {

	// Matches tests whether the given object may be serialized by the current
	// ObjectSerializer.
	//
	// If this method returns true, Serialize will be called on the object and no
	// further ObjectSerializers will be tested.
	Matches(object any) bool

	// Serialize serializes the given object into an io.Reader instance which will
	// be passed to the HTTP client caller.
	Serialize(object any) (io.Reader, error)

	// ContentType returns the content type of the serialized data that this
	// ObjectSerializer returns.
	ContentType() string
}
