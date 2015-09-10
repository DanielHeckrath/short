package client

import (
	"time"

	"github.com/DanielHeckrath/short/pb"
	"github.com/juju/errors"
	"github.com/youtube/vitess/go/pools"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Client is a wrapper for a ShortClient that satisfies the ResourcePools Resource interface
type Client struct {
	pb.ShortClient
	cc *grpc.ClientConn
}

// Close closes the underlying grpc connection
func (pc *Client) Close() {
	pc.cc.Close()
}

// New instantiates a new short grpc client
func New(addr string, options ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(addr, options...)
	if err != nil {
		return nil, errors.Annotate(err, "Unable to connect to user service address")
	}

	psc := pb.NewShortClient(conn)
	client := &Client{
		psc,
		conn,
	}
	return client, nil
}

// Pool is a type safe wrapper around a vitess resource pool
type Pool struct {
	*pools.ResourcePool
}

// Get will return the next available client.
// If capacity has not been reached, it will create a new one using the factory.
// Otherwise, it will wait till the next resource becomes available or a timeout.
// A timeout of 0 is an indefinite wait.
func (p *Pool) Get(ctx context.Context) (pb.ShortClient, error) {
	res, err := p.ResourcePool.Get(ctx)

	if err != nil {
		return nil, err
	}

	client, ok := res.(*Client)

	if !ok {
		return nil, errors.New("Type error during pool.Get")
	}

	return client, nil
}

// Put will return a resource to the pool.
// For every successful Get, a corresponding Put is required.
// If you no longer need a resource, you will need to call Put(nil) instead of
// returning the closed resource.
// The pool will eventually cause a new resource to be created in its place.
func (p *Pool) Put(client pb.ShortClient) error {
	// handle case for Put(nil)
	if client == nil {
		p.ResourcePool.Put(nil)
		return nil
	}

	var pc *Client
	var ok bool
	if pc, ok = client.(*Client); !ok {
		return errors.New("Only instances returned by pool.Get may be put back into the pool")
	}

	p.ResourcePool.Put(pc)
	return nil
}

// AddressFunc is a function that returns a connection address for grpc
type AddressFunc func() string

func createFactoryFunc(addrFunc AddressFunc, timeout time.Duration) pools.Factory {
	return func() (pools.Resource, error) {
		addr := addrFunc()
		// create new client
		return New(addr, grpc.WithBlock(), grpc.WithTimeout(timeout))
	}
}

// NewPool returns a new short service client pool
func NewPool(addr AddressFunc, capacity, maxCap int, connectionTimeout, idleTimeout time.Duration) (*Pool, error) {
	p := pools.NewResourcePool(createFactoryFunc(addr, connectionTimeout), capacity, maxCap, idleTimeout)

	return &Pool{p}, nil
}
