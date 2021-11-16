// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gRPC

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BidAuctionClientFEClient is the client API for BidAuctionClientFE service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BidAuctionClientFEClient interface {
	SendBidRequest(ctx context.Context, in *BidRequest, opts ...grpc.CallOption) (*BidResponse, error)
	SendResultRequest(ctx context.Context, in *ResultRequest, opts ...grpc.CallOption) (*ResultResponse, error)
}

type bidAuctionClientFEClient struct {
	cc grpc.ClientConnInterface
}

func NewBidAuctionClientFEClient(cc grpc.ClientConnInterface) BidAuctionClientFEClient {
	return &bidAuctionClientFEClient{cc}
}

func (c *bidAuctionClientFEClient) SendBidRequest(ctx context.Context, in *BidRequest, opts ...grpc.CallOption) (*BidResponse, error) {
	out := new(BidResponse)
	err := c.cc.Invoke(ctx, "/Proto.BidAuctionClientFE/SendBidRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bidAuctionClientFEClient) SendResultRequest(ctx context.Context, in *ResultRequest, opts ...grpc.CallOption) (*ResultResponse, error) {
	out := new(ResultResponse)
	err := c.cc.Invoke(ctx, "/Proto.BidAuctionClientFE/SendResultRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BidAuctionClientFEServer is the server API for BidAuctionClientFE service.
// All implementations must embed UnimplementedBidAuctionClientFEServer
// for forward compatibility
type BidAuctionClientFEServer interface {
	SendBidRequest(context.Context, *BidRequest) (*BidResponse, error)
	SendResultRequest(context.Context, *ResultRequest) (*ResultResponse, error)
	mustEmbedUnimplementedBidAuctionClientFEServer()
}

// UnimplementedBidAuctionClientFEServer must be embedded to have forward compatible implementations.
type UnimplementedBidAuctionClientFEServer struct {
}

func (UnimplementedBidAuctionClientFEServer) SendBidRequest(context.Context, *BidRequest) (*BidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendBidRequest not implemented")
}
func (UnimplementedBidAuctionClientFEServer) SendResultRequest(context.Context, *ResultRequest) (*ResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendResultRequest not implemented")
}
func (UnimplementedBidAuctionClientFEServer) mustEmbedUnimplementedBidAuctionClientFEServer() {}

// UnsafeBidAuctionClientFEServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BidAuctionClientFEServer will
// result in compilation errors.
type UnsafeBidAuctionClientFEServer interface {
	mustEmbedUnimplementedBidAuctionClientFEServer()
}

func RegisterBidAuctionClientFEServer(s grpc.ServiceRegistrar, srv BidAuctionClientFEServer) {
	s.RegisterService(&BidAuctionClientFE_ServiceDesc, srv)
}

func _BidAuctionClientFE_SendBidRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BidAuctionClientFEServer).SendBidRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Proto.BidAuctionClientFE/SendBidRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BidAuctionClientFEServer).SendBidRequest(ctx, req.(*BidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BidAuctionClientFE_SendResultRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BidAuctionClientFEServer).SendResultRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Proto.BidAuctionClientFE/SendResultRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BidAuctionClientFEServer).SendResultRequest(ctx, req.(*ResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BidAuctionClientFE_ServiceDesc is the grpc.ServiceDesc for BidAuctionClientFE service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BidAuctionClientFE_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Proto.BidAuctionClientFE",
	HandlerType: (*BidAuctionClientFEServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendBidRequest",
			Handler:    _BidAuctionClientFE_SendBidRequest_Handler,
		},
		{
			MethodName: "SendResultRequest",
			Handler:    _BidAuctionClientFE_SendResultRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gRPC/Proto.proto",
}