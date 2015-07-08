package shorten

import (
	"golang.org/x/net/context"

	"github.com/DanielHeckrath/short/pb"
	"github.com/go-kit/kit/endpoint"
)

// NewShortenEndpoint returns a new endpoint for a shorteners shorten function
func NewShortenEndpoint(s Shortener) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		select {
		default:
		case <-ctx.Done():
			return nil, endpoint.ErrContextCanceled
		}

		req, ok := request.(*pb.ShortenRequest)
		if !ok {
			return nil, endpoint.ErrBadCast
		}

		url, err := s.Shorten(ctx, req.Url)

		if err != nil {
			return nil, err
		}

		return &pb.ShortenResponse{Url: url}, nil
	}
}

// NewResolveEndpoint returns a new endpoint for a shorteners shorten function
func NewResolveEndpoint(s Shortener) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		select {
		default:
		case <-ctx.Done():
			return nil, endpoint.ErrContextCanceled
		}

		req, ok := request.(*pb.ResolveRequest)
		if !ok {
			return nil, endpoint.ErrBadCast
		}

		url, err := s.Resolve(ctx, req.Key)

		if err != nil {
			return nil, err
		}

		return &pb.ResolveResponse{Url: url}, nil
	}
}

// NewInfoEndpoint returns a new endpoint for a shorteners info function
func NewInfoEndpoint(s Shortener) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		select {
		default:
		case <-ctx.Done():
			return nil, endpoint.ErrContextCanceled
		}

		req, ok := request.(*pb.InfoRequest)
		if !ok {
			return nil, endpoint.ErrBadCast
		}

		url, err := s.Info(ctx, req.Key)

		if err != nil {
			return nil, err
		}

		return &pb.InfoResponse{Url: url}, nil
	}
}

// NewLatestEndpoint returns a new endpoint for a shorteners shorten function
func NewLatestEndpoint(s Shortener) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		select {
		default:
		case <-ctx.Done():
			return nil, endpoint.ErrContextCanceled
		}

		req, ok := request.(*pb.LatestRequest)
		if !ok {
			return nil, endpoint.ErrBadCast
		}

		urls, err := s.Latest(ctx, req.Count)

		if err != nil {
			return nil, err
		}

		return &pb.LatestResponse{Urls: urls}, nil
	}
}
