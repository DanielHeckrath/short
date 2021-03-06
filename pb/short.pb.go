// Code generated by protoc-gen-go.
// source: short.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	short.proto

It has these top-level messages:
	ShortURL
	ShortenRequest
	ShortenResponse
	ResolveRequest
	ResolveResponse
	InfoRequest
	InfoResponse
	LatestRequest
	LatestResponse
*/
package pb

import proto "github.com/golang/protobuf/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type ShortURL struct {
	Key          string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	ShortUrl     string `protobuf:"bytes,2,opt,name=short_url" json:"short_url,omitempty"`
	LongUrl      string `protobuf:"bytes,3,opt,name=long_url" json:"long_url,omitempty"`
	CreationDate int64  `protobuf:"varint,4,opt,name=creation_date" json:"creation_date,omitempty"`
	Clicks       int64  `protobuf:"varint,5,opt,name=clicks" json:"clicks,omitempty"`
}

func (m *ShortURL) Reset()         { *m = ShortURL{} }
func (m *ShortURL) String() string { return proto.CompactTextString(m) }
func (*ShortURL) ProtoMessage()    {}

type ShortenRequest struct {
	Url string `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
}

func (m *ShortenRequest) Reset()         { *m = ShortenRequest{} }
func (m *ShortenRequest) String() string { return proto.CompactTextString(m) }
func (*ShortenRequest) ProtoMessage()    {}

type ShortenResponse struct {
	Url *ShortURL `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
}

func (m *ShortenResponse) Reset()         { *m = ShortenResponse{} }
func (m *ShortenResponse) String() string { return proto.CompactTextString(m) }
func (*ShortenResponse) ProtoMessage()    {}

func (m *ShortenResponse) GetUrl() *ShortURL {
	if m != nil {
		return m.Url
	}
	return nil
}

type ResolveRequest struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *ResolveRequest) Reset()         { *m = ResolveRequest{} }
func (m *ResolveRequest) String() string { return proto.CompactTextString(m) }
func (*ResolveRequest) ProtoMessage()    {}

type ResolveResponse struct {
	Url *ShortURL `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
}

func (m *ResolveResponse) Reset()         { *m = ResolveResponse{} }
func (m *ResolveResponse) String() string { return proto.CompactTextString(m) }
func (*ResolveResponse) ProtoMessage()    {}

func (m *ResolveResponse) GetUrl() *ShortURL {
	if m != nil {
		return m.Url
	}
	return nil
}

type InfoRequest struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *InfoRequest) Reset()         { *m = InfoRequest{} }
func (m *InfoRequest) String() string { return proto.CompactTextString(m) }
func (*InfoRequest) ProtoMessage()    {}

type InfoResponse struct {
	Url *ShortURL `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
}

func (m *InfoResponse) Reset()         { *m = InfoResponse{} }
func (m *InfoResponse) String() string { return proto.CompactTextString(m) }
func (*InfoResponse) ProtoMessage()    {}

func (m *InfoResponse) GetUrl() *ShortURL {
	if m != nil {
		return m.Url
	}
	return nil
}

type LatestRequest struct {
	Count int64 `protobuf:"varint,1,opt,name=count" json:"count,omitempty"`
}

func (m *LatestRequest) Reset()         { *m = LatestRequest{} }
func (m *LatestRequest) String() string { return proto.CompactTextString(m) }
func (*LatestRequest) ProtoMessage()    {}

type LatestResponse struct {
	Urls []*ShortURL `protobuf:"bytes,1,rep,name=urls" json:"urls,omitempty"`
}

func (m *LatestResponse) Reset()         { *m = LatestResponse{} }
func (m *LatestResponse) String() string { return proto.CompactTextString(m) }
func (*LatestResponse) ProtoMessage()    {}

func (m *LatestResponse) GetUrls() []*ShortURL {
	if m != nil {
		return m.Urls
	}
	return nil
}

func init() {
}

// Client API for Short service

type ShortClient interface {
	Shorten(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error)
	Resolve(ctx context.Context, in *ResolveRequest, opts ...grpc.CallOption) (*ResolveResponse, error)
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error)
	Latest(ctx context.Context, in *LatestRequest, opts ...grpc.CallOption) (*LatestResponse, error)
}

type shortClient struct {
	cc *grpc.ClientConn
}

func NewShortClient(cc *grpc.ClientConn) ShortClient {
	return &shortClient{cc}
}

func (c *shortClient) Shorten(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error) {
	out := new(ShortenResponse)
	err := grpc.Invoke(ctx, "/pb.Short/Shorten", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortClient) Resolve(ctx context.Context, in *ResolveRequest, opts ...grpc.CallOption) (*ResolveResponse, error) {
	out := new(ResolveResponse)
	err := grpc.Invoke(ctx, "/pb.Short/Resolve", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	out := new(InfoResponse)
	err := grpc.Invoke(ctx, "/pb.Short/Info", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortClient) Latest(ctx context.Context, in *LatestRequest, opts ...grpc.CallOption) (*LatestResponse, error) {
	out := new(LatestResponse)
	err := grpc.Invoke(ctx, "/pb.Short/Latest", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Short service

type ShortServer interface {
	Shorten(context.Context, *ShortenRequest) (*ShortenResponse, error)
	Resolve(context.Context, *ResolveRequest) (*ResolveResponse, error)
	Info(context.Context, *InfoRequest) (*InfoResponse, error)
	Latest(context.Context, *LatestRequest) (*LatestResponse, error)
}

func RegisterShortServer(s *grpc.Server, srv ShortServer) {
	s.RegisterService(&_Short_serviceDesc, srv)
}

func _Short_Shorten_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(ShortenRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(ShortServer).Shorten(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Short_Resolve_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(ResolveRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(ShortServer).Resolve(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Short_Info_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(InfoRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(ShortServer).Info(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Short_Latest_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(LatestRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(ShortServer).Latest(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Short_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Short",
	HandlerType: (*ShortServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Shorten",
			Handler:    _Short_Shorten_Handler,
		},
		{
			MethodName: "Resolve",
			Handler:    _Short_Resolve_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _Short_Info_Handler,
		},
		{
			MethodName: "Latest",
			Handler:    _Short_Latest_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
