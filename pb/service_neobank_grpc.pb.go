// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: service_neobank.proto

package pb

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

// NeoBankClient is the client API for NeoBank service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NeoBankClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error)
}

type neoBankClient struct {
	cc grpc.ClientConnInterface
}

func NewNeoBankClient(cc grpc.ClientConnInterface) NeoBankClient {
	return &neoBankClient{cc}
}

func (c *neoBankClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/pb.NeoBank/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *neoBankClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error) {
	out := new(LoginUserResponse)
	err := c.cc.Invoke(ctx, "/pb.NeoBank/LoginUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NeoBankServer is the server API for NeoBank service.
// All implementations must embed UnimplementedNeoBankServer
// for forward compatibility
type NeoBankServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error)
	mustEmbedUnimplementedNeoBankServer()
}

// UnimplementedNeoBankServer must be embedded to have forward compatible implementations.
type UnimplementedNeoBankServer struct {
}

func (UnimplementedNeoBankServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedNeoBankServer) LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedNeoBankServer) mustEmbedUnimplementedNeoBankServer() {}

// UnsafeNeoBankServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NeoBankServer will
// result in compilation errors.
type UnsafeNeoBankServer interface {
	mustEmbedUnimplementedNeoBankServer()
}

func RegisterNeoBankServer(s grpc.ServiceRegistrar, srv NeoBankServer) {
	s.RegisterService(&NeoBank_ServiceDesc, srv)
}

func _NeoBank_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeoBankServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.NeoBank/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeoBankServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NeoBank_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NeoBankServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.NeoBank/LoginUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NeoBankServer).LoginUser(ctx, req.(*LoginUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NeoBank_ServiceDesc is the grpc.ServiceDesc for NeoBank service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NeoBank_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.NeoBank",
	HandlerType: (*NeoBankServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _NeoBank_CreateUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _NeoBank_LoginUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_neobank.proto",
}
