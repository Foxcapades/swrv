package swrv

import (
	"fmt"
	"io"
	"strings"
)

type defaultObjectSerializer struct{}

func (d defaultObjectSerializer) Matches(object any) bool {
	return true
}

func (d defaultObjectSerializer) Serialize(object any) (io.Reader, error) {
	return strings.NewReader(fmt.Sprint(object)), nil
}

func (d defaultObjectSerializer) ContentType() string {
	return ContentTypeTextPlain
}
