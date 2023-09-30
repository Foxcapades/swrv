package swrv

// RequestContext is a map of arbitrary state that may be attached to a Request
// instance as it passes through the various stages of the request handling
// process.
type RequestContext interface {

	// Has tests whether the RequestContext contains and entry with the given key.
	Has(key string) bool

	// Get returns the value stored at the given key.
	Get(key string) any

	// Put sets the given value into the RequestContext at the given key.
	Put(key string, val any)

	// Len returns the current size of the RequestContext instance.
	Len() int

	// IsEmpty tests whether this RequestContext has a length of 0.
	IsEmpty() bool
}

type requestContext map[string]interface{}

func (r requestContext) Has(key string) bool {
	_, ok := r[key]
	return ok
}

func (r requestContext) Get(key string) any {
	return r[key]
}

func (r requestContext) Put(key string, val interface{}) {
	r[key] = val
}

func (r requestContext) Len() int {
	return len(r)
}

func (r requestContext) IsEmpty() bool {
	return len(r) == 0
}
