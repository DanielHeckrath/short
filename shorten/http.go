package shorten

import (
	"encoding/json"
	stdhttp "net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	transport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/juju/errors"
	"golang.org/x/net/context"

	"github.com/DanielHeckrath/short/http"
	"github.com/DanielHeckrath/short/pb"
)

func encodeShortenResponse(w stdhttp.ResponseWriter, response interface{}) error {
	shorten, ok := response.(*pb.ShortenResponse)

	if !ok {
		return errors.New("Response must be a ShortenResponse")
	}

	return json.NewEncoder(w).Encode(shorten.Url)
}

func decodeShortenRequest(r *stdhttp.Request) (interface{}, error) {
	url := r.FormValue("url")

	if url == "" {
		return nil, errors.New("URL can't be empty")
	}

	return &pb.ShortenRequest{Url: url}, nil
}

func encodeResolveResponse(w stdhttp.ResponseWriter, response interface{}) error {
	resolve, ok := response.(*pb.ResolveResponse)

	if !ok {
		return errors.New("Response must be a ResolveResponse")
	}

	return json.NewEncoder(w).Encode(resolve.Url)
}

func decodeResolveRequest(r *stdhttp.Request) (interface{}, error) {
	key := mux.Vars(r)["key"]

	if key == "" {
		return nil, errors.New("Key can't be empty")
	}

	return &pb.ResolveRequest{Key: key}, nil
}

func encodeInfoResponse(w stdhttp.ResponseWriter, response interface{}) error {
	info, ok := response.(*pb.InfoResponse)

	if !ok {
		return errors.New("Response must be a InfoResponse")
	}

	return json.NewEncoder(w).Encode(info.Url)
}

func decodeInfoRequest(r *stdhttp.Request) (interface{}, error) {
	key := mux.Vars(r)["key"]

	if key == "" {
		return nil, errors.New("Key can't be empty")
	}

	return &pb.InfoRequest{Key: key}, nil
}

func encodeLatestResponse(w stdhttp.ResponseWriter, response interface{}) error {
	latest, ok := response.(*pb.LatestResponse)

	if !ok {
		return errors.New("Response must be a LatestResponse")
	}

	return json.NewEncoder(w).Encode(latest.Urls)
}

func decodeLatestRequest(r *stdhttp.Request) (interface{}, error) {
	data := mux.Vars(r)["count"]
	count, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		count = 10
	}

	return &pb.LatestRequest{Count: count}, nil
}

// NewShortenHandler returns a http.Handler for a shorten function endpoint
func NewShortenHandler(ctx context.Context, endpoint endpoint.Endpoint, before []transport.BeforeFunc, after []transport.AfterFunc) stdhttp.Handler {
	return http.NewHandler(ctx, endpoint, decodeShortenRequest, encodeShortenResponse, before, after)
}

type resolveHandler struct {
	ctx      context.Context
	e        endpoint.Endpoint
	notFound string
}

func (h resolveHandler) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	req, err := decodeResolveRequest(r)

	if err != nil {
		stdhttp.Error(w, "Unable to decode request", stdhttp.StatusBadRequest)
		return
	}

	resolve, ok := req.(*pb.ResolveRequest)

	if !ok {
		stdhttp.Error(w, "Unable to decode request", stdhttp.StatusBadRequest)
		return
	}

	res, err := h.e(h.ctx, resolve)

	if err != nil {
		stdhttp.Redirect(w, r, h.notFound, stdhttp.StatusMovedPermanently)
		return
	}

	response, ok := res.(*pb.ResolveResponse)

	if !ok {
		stdhttp.Error(w, "Response must be a ResolveResponse", stdhttp.StatusInternalServerError)
		return
	}

	stdhttp.Redirect(w, r, response.Url.LongUrl, stdhttp.StatusMovedPermanently)
}

// NewResolveHandler returns a http.Handler for a resolve function endpoint
func NewResolveHandler(ctx context.Context, endpoint endpoint.Endpoint, notFound string) stdhttp.Handler {
	return resolveHandler{
		ctx:      ctx,
		e:        endpoint,
		notFound: notFound,
	}
}

// NewInfoHandler returns a http.Handler for a resolve function endpoint
func NewInfoHandler(ctx context.Context, endpoint endpoint.Endpoint, before []transport.BeforeFunc, after []transport.AfterFunc) stdhttp.Handler {
	return http.NewHandler(ctx, endpoint, decodeInfoRequest, encodeInfoResponse, before, after)
}

// NewLatestHandler returns a http.Handler for a latest function endpoint
func NewLatestHandler(ctx context.Context, endpoint endpoint.Endpoint, before []transport.BeforeFunc, after []transport.AfterFunc) stdhttp.Handler {
	return http.NewHandler(ctx, endpoint, decodeLatestRequest, encodeLatestResponse, before, after)
}
