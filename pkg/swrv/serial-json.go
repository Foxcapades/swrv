package swrv

import (
	"bytes"
	"encoding/json"
	"io"
)

func defaultMatcher(_ any) bool {
	return true
}

// NewDefaultJSONObjectSerializer returns an ObjectSerializer instance that will
// match ALL objects and attempt to serialize them as JSON.
func NewDefaultJSONObjectSerializer() ObjectSerializer {
	return NewJSONObjectSerializer(defaultMatcher)
}

// NewJSONObjectSerializer returns an ObjectSerializer instance that will match
// only the objects that the given MatcherFn instance returns true for, and will
// attempt to serialize them as JSON.
func NewJSONObjectSerializer(fn MatcherFn) ObjectSerializer {
	return jsonObjectSerializer{fn}
}

type jsonObjectSerializer struct {
	matcher MatcherFn
}

func (j jsonObjectSerializer) Matches(object any) bool {
	return j.matcher(object)
}

func (j jsonObjectSerializer) Serialize(object any) (io.Reader, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	return buffer, encoder.Encode(object)
}

func (j jsonObjectSerializer) ContentType() string {
	return ContentTypeApplicationJSON
}
