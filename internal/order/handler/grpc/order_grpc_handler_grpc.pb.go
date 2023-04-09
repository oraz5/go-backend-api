// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: internal/order/handler/grpc/order_grpc_handler.proto

package grpc_server

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

// OrderHandlerClient is the client API for OrderHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderHandlerClient interface {
	OrderList(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
}

type orderHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderHandlerClient(cc grpc.ClientConnInterface) OrderHandlerClient {
	return &orderHandlerClient{cc}
}

func (c *orderHandlerClient) OrderList(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, "/grpc_handler.OrderHandler/OrderList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderHandlerServer is the server API for OrderHandler service.
// All implementations must embed UnimplementedOrderHandlerServer
// for forward compatibility
type OrderHandlerServer interface {
	OrderList(context.Context, *OrderRequest) (*OrderResponse, error)
	mustEmbedUnimplementedOrderHandlerServer()
}

// UnimplementedOrderHandlerServer must be embedded to have forward compatible implementations.
type UnimplementedOrderHandlerServer struct {
}

func (UnimplementedOrderHandlerServer) OrderList(context.Context, *OrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderList not implemented")
}
func (UnimplementedOrderHandlerServer) mustEmbedUnimplementedOrderHandlerServer() {}

// UnsafeOrderHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderHandlerServer will
// result in compilation errors.
type UnsafeOrderHandlerServer interface {
	mustEmbedUnimplementedOrderHandlerServer()
}

func RegisterOrderHandlerServer(s grpc.ServiceRegistrar, srv OrderHandlerServer) {
	s.RegisterService(&OrderHandler_ServiceDesc, srv)
}

func _OrderHandler_OrderList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderHandlerServer).OrderList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_handler.OrderHandler/OrderList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderHandlerServer).OrderList(ctx, req.(*OrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderHandler_ServiceDesc is the grpc.ServiceDesc for OrderHandler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderHandler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc_handler.OrderHandler",
	HandlerType: (*OrderHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OrderList",
			Handler:    _OrderHandler_OrderList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/order/handler/grpc/order_grpc_handler.proto",
}