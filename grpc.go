package main

import (
	"golang.org/x/net/context"

	"github.com/DanielHeckrath/short/pb"
	"github.com/go-kit/kit/endpoint"
)

// A binding wraps an Endpoint so that it's usable by a transport. grpcBinding
// makes an Endpoint usable over gRPC.
type grpcBinding struct {
	shorten endpoint.Endpoint
	resolve endpoint.Endpoint
	latest  endpoint.Endpoint
}

func (b grpcBinding) Resolve(ctx0 context.Context, req *pb.ResolveRequest) (*pb.ResolveResponse, error) {
	var (
		ctx, cancel = context.WithCancel(ctx0)
		errs        = make(chan error, 1)
		replies     = make(chan *pb.ResolveResponse, 1)
	)
	defer cancel()
	go func() {
		r, err := b.latest(ctx, req)
		if err != nil {
			errs <- err
			return
		}
		resp, ok := r.(*pb.ResolveResponse)
		if !ok {
			errs <- endpoint.ErrBadCast
			return
		}
		replies <- resp
	}()
	select {
	case <-ctx.Done():
		return nil, context.DeadlineExceeded
	case err := <-errs:
		return nil, err
	case reply := <-replies:
		return reply, nil
	}
}

func (b grpcBinding) Latest(ctx0 context.Context, req *pb.LatestRequest) (*pb.LatestResponse, error) {
	var (
		ctx, cancel = context.WithCancel(ctx0)
		errs        = make(chan error, 1)
		replies     = make(chan *pb.LatestResponse, 1)
	)
	defer cancel()
	go func() {
		r, err := b.latest(ctx, req)
		if err != nil {
			errs <- err
			return
		}
		resp, ok := r.(*pb.LatestResponse)
		if !ok {
			errs <- endpoint.ErrBadCast
			return
		}
		replies <- resp
	}()
	select {
	case <-ctx.Done():
		return nil, context.DeadlineExceeded
	case err := <-errs:
		return nil, err
	case reply := <-replies:
		return reply, nil
	}
}

// Shorten implements the proto3 Short service by forwarding to the wrapped Endpoint.
func (b grpcBinding) Shorten(ctx0 context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	var (
		ctx, cancel = context.WithCancel(ctx0)
		errs        = make(chan error, 1)
		replies     = make(chan *pb.ShortenResponse, 1)
	)
	defer cancel()
	go func() {
		r, err := b.shorten(ctx, req)
		if err != nil {
			errs <- err
			return
		}
		resp, ok := r.(*pb.ShortenResponse)
		if !ok {
			errs <- endpoint.ErrBadCast
			return
		}
		replies <- resp
	}()
	select {
	case <-ctx.Done():
		return nil, context.DeadlineExceeded
	case err := <-errs:
		return nil, err
	case reply := <-replies:
		return reply, nil
	}
}
