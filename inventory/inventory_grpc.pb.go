// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: inventory/inventory.proto

package inventory

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Inventory_GetUserInventory_FullMethodName = "/inventory.Inventory/GetUserInventory"
	Inventory_GetIngredient_FullMethodName    = "/inventory.Inventory/GetIngredient"
	Inventory_CreateIngredient_FullMethodName = "/inventory.Inventory/CreateIngredient"
	Inventory_UpdateIngredient_FullMethodName = "/inventory.Inventory/UpdateIngredient"
	Inventory_DeleteIngredient_FullMethodName = "/inventory.Inventory/DeleteIngredient"
)

// InventoryClient is the client API for Inventory service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The Inventory service definition.
type InventoryClient interface {
	GetUserInventory(ctx context.Context, in *GetInventoryRequest, opts ...grpc.CallOption) (*GetUserInventoryResponse, error)
	GetIngredient(ctx context.Context, in *GetIngredientRequest, opts ...grpc.CallOption) (*IngredientResponse, error)
	CreateIngredient(ctx context.Context, in *PostIngredientRequest, opts ...grpc.CallOption) (*IngredientResponse, error)
	UpdateIngredient(ctx context.Context, in *PostIngredientRequest, opts ...grpc.CallOption) (*IngredientResponse, error)
	DeleteIngredient(ctx context.Context, in *DeleteIngredientRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type inventoryClient struct {
	cc grpc.ClientConnInterface
}

func NewInventoryClient(cc grpc.ClientConnInterface) InventoryClient {
	return &inventoryClient{cc}
}

func (c *inventoryClient) GetUserInventory(ctx context.Context, in *GetInventoryRequest, opts ...grpc.CallOption) (*GetUserInventoryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserInventoryResponse)
	err := c.cc.Invoke(ctx, Inventory_GetUserInventory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryClient) GetIngredient(ctx context.Context, in *GetIngredientRequest, opts ...grpc.CallOption) (*IngredientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IngredientResponse)
	err := c.cc.Invoke(ctx, Inventory_GetIngredient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryClient) CreateIngredient(ctx context.Context, in *PostIngredientRequest, opts ...grpc.CallOption) (*IngredientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IngredientResponse)
	err := c.cc.Invoke(ctx, Inventory_CreateIngredient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryClient) UpdateIngredient(ctx context.Context, in *PostIngredientRequest, opts ...grpc.CallOption) (*IngredientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IngredientResponse)
	err := c.cc.Invoke(ctx, Inventory_UpdateIngredient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryClient) DeleteIngredient(ctx context.Context, in *DeleteIngredientRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Inventory_DeleteIngredient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InventoryServer is the server API for Inventory service.
// All implementations must embed UnimplementedInventoryServer
// for forward compatibility.
//
// The Inventory service definition.
type InventoryServer interface {
	GetUserInventory(context.Context, *GetInventoryRequest) (*GetUserInventoryResponse, error)
	GetIngredient(context.Context, *GetIngredientRequest) (*IngredientResponse, error)
	CreateIngredient(context.Context, *PostIngredientRequest) (*IngredientResponse, error)
	UpdateIngredient(context.Context, *PostIngredientRequest) (*IngredientResponse, error)
	DeleteIngredient(context.Context, *DeleteIngredientRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedInventoryServer()
}

// UnimplementedInventoryServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedInventoryServer struct{}

func (UnimplementedInventoryServer) GetUserInventory(context.Context, *GetInventoryRequest) (*GetUserInventoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInventory not implemented")
}
func (UnimplementedInventoryServer) GetIngredient(context.Context, *GetIngredientRequest) (*IngredientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIngredient not implemented")
}
func (UnimplementedInventoryServer) CreateIngredient(context.Context, *PostIngredientRequest) (*IngredientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateIngredient not implemented")
}
func (UnimplementedInventoryServer) UpdateIngredient(context.Context, *PostIngredientRequest) (*IngredientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateIngredient not implemented")
}
func (UnimplementedInventoryServer) DeleteIngredient(context.Context, *DeleteIngredientRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteIngredient not implemented")
}
func (UnimplementedInventoryServer) mustEmbedUnimplementedInventoryServer() {}
func (UnimplementedInventoryServer) testEmbeddedByValue()                   {}

// UnsafeInventoryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InventoryServer will
// result in compilation errors.
type UnsafeInventoryServer interface {
	mustEmbedUnimplementedInventoryServer()
}

func RegisterInventoryServer(s grpc.ServiceRegistrar, srv InventoryServer) {
	// If the following call pancis, it indicates UnimplementedInventoryServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Inventory_ServiceDesc, srv)
}

func _Inventory_GetUserInventory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInventoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServer).GetUserInventory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Inventory_GetUserInventory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServer).GetUserInventory(ctx, req.(*GetInventoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Inventory_GetIngredient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIngredientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServer).GetIngredient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Inventory_GetIngredient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServer).GetIngredient(ctx, req.(*GetIngredientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Inventory_CreateIngredient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostIngredientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServer).CreateIngredient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Inventory_CreateIngredient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServer).CreateIngredient(ctx, req.(*PostIngredientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Inventory_UpdateIngredient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostIngredientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServer).UpdateIngredient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Inventory_UpdateIngredient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServer).UpdateIngredient(ctx, req.(*PostIngredientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Inventory_DeleteIngredient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteIngredientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServer).DeleteIngredient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Inventory_DeleteIngredient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServer).DeleteIngredient(ctx, req.(*DeleteIngredientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Inventory_ServiceDesc is the grpc.ServiceDesc for Inventory service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Inventory_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "inventory.Inventory",
	HandlerType: (*InventoryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserInventory",
			Handler:    _Inventory_GetUserInventory_Handler,
		},
		{
			MethodName: "GetIngredient",
			Handler:    _Inventory_GetIngredient_Handler,
		},
		{
			MethodName: "CreateIngredient",
			Handler:    _Inventory_CreateIngredient_Handler,
		},
		{
			MethodName: "UpdateIngredient",
			Handler:    _Inventory_UpdateIngredient_Handler,
		},
		{
			MethodName: "DeleteIngredient",
			Handler:    _Inventory_DeleteIngredient_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "inventory/inventory.proto",
}