package swrv

import "strings"

func newEmptyResponseError(error string) Response {
	return NewResponse().WithCode(500).WithBody(strings.NewReader(error))
}
