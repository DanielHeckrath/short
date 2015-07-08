package http

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	transport "github.com/go-kit/kit/transport/http"
)

// NewHandler return a http.Handler for a function endpoint
//
// NewHandler will add an AfterFunc that sets the response content type to application/json; charset=utf-8
func NewHandler(ctx context.Context,
	endpoint endpoint.Endpoint,
	decode transport.DecodeFunc,
	encode transport.EncodeFunc,
	before []transport.BeforeFunc,
	after []transport.AfterFunc,
) http.Handler {
	return transport.Server{
		Context:    ctx,
		Endpoint:   endpoint,
		DecodeFunc: decode,
		EncodeFunc: encode,
		Before:     before,
		After: append([]transport.AfterFunc{
			transport.SetContentType("application/json; charset=utf-8"),
		}, after...),
	}
}
