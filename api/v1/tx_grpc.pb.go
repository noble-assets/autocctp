// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: noble/autocctp/v1/tx.proto

package autocctpv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Msg_RegisterAccount_FullMethodName             = "/noble.autocctp.v1.Msg/RegisterAccount"
	Msg_RegisterAccountSignerlessly_FullMethodName = "/noble.autocctp.v1.Msg/RegisterAccountSignerlessly"
	Msg_ClearAccount_FullMethodName                = "/noble.autocctp.v1.Msg/ClearAccount"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	RegisterAccount(ctx context.Context, in *MsgRegisterAccount, opts ...grpc.CallOption) (*MsgRegisterAccountResponse, error)
	RegisterAccountSignerlessly(ctx context.Context, in *MsgRegisterAccountSignerlessly, opts ...grpc.CallOption) (*MsgRegisterAccountSignerlesslyResponse, error)
	ClearAccount(ctx context.Context, in *MsgClearAccount, opts ...grpc.CallOption) (*MsgClearAccountResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) RegisterAccount(ctx context.Context, in *MsgRegisterAccount, opts ...grpc.CallOption) (*MsgRegisterAccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgRegisterAccountResponse)
	err := c.cc.Invoke(ctx, Msg_RegisterAccount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RegisterAccountSignerlessly(ctx context.Context, in *MsgRegisterAccountSignerlessly, opts ...grpc.CallOption) (*MsgRegisterAccountSignerlesslyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgRegisterAccountSignerlesslyResponse)
	err := c.cc.Invoke(ctx, Msg_RegisterAccountSignerlessly_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ClearAccount(ctx context.Context, in *MsgClearAccount, opts ...grpc.CallOption) (*MsgClearAccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgClearAccountResponse)
	err := c.cc.Invoke(ctx, Msg_ClearAccount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility.
type MsgServer interface {
	RegisterAccount(context.Context, *MsgRegisterAccount) (*MsgRegisterAccountResponse, error)
	RegisterAccountSignerlessly(context.Context, *MsgRegisterAccountSignerlessly) (*MsgRegisterAccountSignerlesslyResponse, error)
	ClearAccount(context.Context, *MsgClearAccount) (*MsgClearAccountResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMsgServer struct{}

func (UnimplementedMsgServer) RegisterAccount(context.Context, *MsgRegisterAccount) (*MsgRegisterAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterAccount not implemented")
}
func (UnimplementedMsgServer) RegisterAccountSignerlessly(context.Context, *MsgRegisterAccountSignerlessly) (*MsgRegisterAccountSignerlesslyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterAccountSignerlessly not implemented")
}
func (UnimplementedMsgServer) ClearAccount(context.Context, *MsgClearAccount) (*MsgClearAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearAccount not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}
func (UnimplementedMsgServer) testEmbeddedByValue()             {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	// If the following call pancis, it indicates UnimplementedMsgServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_RegisterAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterAccount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RegisterAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterAccount(ctx, req.(*MsgRegisterAccount))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RegisterAccountSignerlessly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterAccountSignerlessly)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterAccountSignerlessly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RegisterAccountSignerlessly_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterAccountSignerlessly(ctx, req.(*MsgRegisterAccountSignerlessly))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClearAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClearAccount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClearAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_ClearAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClearAccount(ctx, req.(*MsgClearAccount))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "noble.autocctp.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterAccount",
			Handler:    _Msg_RegisterAccount_Handler,
		},
		{
			MethodName: "RegisterAccountSignerlessly",
			Handler:    _Msg_RegisterAccountSignerlessly_Handler,
		},
		{
			MethodName: "ClearAccount",
			Handler:    _Msg_ClearAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "noble/autocctp/v1/tx.proto",
}
