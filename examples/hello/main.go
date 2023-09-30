package main

import (
	"fmt"

	"github.com/foxcapades/swrv/pkg/swrv"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	swrv.NewServer("0.0.0.0", 8080).
		WithObjectSerializers(swrv.NewDefaultJSONObjectSerializer()).
		WithControllers(
			swrv.NewController("/hello/{name}", swrv.RequestHandlerFunc(func(request swrv.Request) swrv.Response {
				return swrv.NewResponse().
					WithBody(Response{fmt.Sprintf("Hello %s!", request.URIParam("name"))})
			})),
		).
		Start(nil)
}
