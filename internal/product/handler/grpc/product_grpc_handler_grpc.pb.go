// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: internal/product/handler/grpc/product_grpc_handler.proto

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

// ProductHandlerClient is the client API for ProductHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProductHandlerClient interface {
	Products(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*ProductReponse, error)
	ProductSingle(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*ProductReponse, error)
}

type productHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewProductHandlerClient(cc grpc.ClientConnInterface) ProductHandlerClient {
	return &productHandlerClient{cc}
}

func (c *productHandlerClient) Products(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*ProductReponse, error) {
	out := new(ProductReponse)
	err := c.cc.Invoke(ctx, "/grpc_handler.ProductHandler/Products", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productHandlerClient) ProductSingle(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*ProductReponse, error) {
	out := new(ProductReponse)
	err := c.cc.Invoke(ctx, "/grpc_handler.ProductHandler/ProductSingle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductHandlerServer is the server API for ProductHandler service.
// All implementations must embed UnimplementedProductHandlerServer
// for forward compatibility
type ProductHandlerServer interface {
	Products(context.Context, *ProductRequest) (*ProductReponse, error)
	ProductSingle(context.Context, *ProductRequest) (*ProductReponse, error)
	mustEmbedUnimplementedProductHandlerServer()
}

// UnimplementedProductHandlerServer must be embedded to have forward compatible implementations.
type UnimplementedProductHandlerServer struct {
}

func (UnimplementedProductHandlerServer) Products(context.Context, *ProductRequest) (*ProductReponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Products not implemented")
}
func (UnimplementedProductHandlerServer) ProductSingle(context.Context, *ProductRequest) (*ProductReponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProductSingle not implemented")
}
func (UnimplementedProductHandlerServer) mustEmbedUnimplementedProductHandlerServer() {}

// UnsafeProductHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProductHandlerServer will
// result in compilation errors.
type UnsafeProductHandlerServer interface {
	mustEmbedUnimplementedProductHandlerServer()
}

func RegisterProductHandlerServer(s grpc.ServiceRegistrar, srv ProductHandlerServer) {
	s.RegisterService(&ProductHandler_ServiceDesc, srv)
}

func _ProductHandler_Products_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductHandlerServer).Products(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_handler.ProductHandler/Products",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductHandlerServer).Products(ctx, req.(*ProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductHandler_ProductSingle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductHandlerServer).ProductSingle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_handler.ProductHandler/ProductSingle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductHandlerServer).ProductSingle(ctx, req.(*ProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProductHandler_ServiceDesc is the grpc.ServiceDesc for ProductHandler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProductHandler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc_handler.ProductHandler",
	HandlerType: (*ProductHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Products",
			Handler:    _ProductHandler_Products_Handler,
		},
		{
			MethodName: "ProductSingle",
			Handler:    _ProductHandler_ProductSingle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/product/handler/grpc/product_grpc_handler.proto",
}
