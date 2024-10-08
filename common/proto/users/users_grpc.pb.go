// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.2
// source: common/proto/users/users.proto

package proto

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

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	RegisterUser(ctx context.Context, in *UserRegistrationPayload, opts ...grpc.CallOption) (*RegisterMessage, error)
	LoginUser(ctx context.Context, in *UserLoginPayload, opts ...grpc.CallOption) (*LoginResponse, error)
	CreateRating(ctx context.Context, in *ReviewRequest, opts ...grpc.CallOption) (*ReviewMessage, error)
	GetUserByID(ctx context.Context, in *UserIDMessage, opts ...grpc.CallOption) (*User, error)
	ChargeUserWallet(ctx context.Context, in *ChargeRequest, opts ...grpc.CallOption) (*ChargeResponse, error)
	CreditUserWallet(ctx context.Context, in *ChargeRequest, opts ...grpc.CallOption) (*ChargeResponse, error)
	VerifyOTP(ctx context.Context, in *OTPVerifyPayload, opts ...grpc.CallOption) (*OTPVerifyResponse, error)
	ResendOTP(ctx context.Context, in *OTPResendPayload, opts ...grpc.CallOption) (*OTPResendResponse, error)
	GetWalletBalance(ctx context.Context, in *WalletBalanceRequest, opts ...grpc.CallOption) (*WalletBalanceResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) RegisterUser(ctx context.Context, in *UserRegistrationPayload, opts ...grpc.CallOption) (*RegisterMessage, error) {
	out := new(RegisterMessage)
	err := c.cc.Invoke(ctx, "/proto.UserService/RegisterUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) LoginUser(ctx context.Context, in *UserLoginPayload, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/proto.UserService/LoginUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateRating(ctx context.Context, in *ReviewRequest, opts ...grpc.CallOption) (*ReviewMessage, error) {
	out := new(ReviewMessage)
	err := c.cc.Invoke(ctx, "/proto.UserService/CreateRating", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserByID(ctx context.Context, in *UserIDMessage, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/proto.UserService/GetUserByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ChargeUserWallet(ctx context.Context, in *ChargeRequest, opts ...grpc.CallOption) (*ChargeResponse, error) {
	out := new(ChargeResponse)
	err := c.cc.Invoke(ctx, "/proto.UserService/ChargeUserWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreditUserWallet(ctx context.Context, in *ChargeRequest, opts ...grpc.CallOption) (*ChargeResponse, error) {
	out := new(ChargeResponse)
	err := c.cc.Invoke(ctx, "/proto.UserService/CreditUserWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) VerifyOTP(ctx context.Context, in *OTPVerifyPayload, opts ...grpc.CallOption) (*OTPVerifyResponse, error) {
	out := new(OTPVerifyResponse)
	err := c.cc.Invoke(ctx, "/proto.UserService/VerifyOTP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ResendOTP(ctx context.Context, in *OTPResendPayload, opts ...grpc.CallOption) (*OTPResendResponse, error) {
	out := new(OTPResendResponse)
	err := c.cc.Invoke(ctx, "/proto.UserService/ResendOTP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetWalletBalance(ctx context.Context, in *WalletBalanceRequest, opts ...grpc.CallOption) (*WalletBalanceResponse, error) {
	out := new(WalletBalanceResponse)
	err := c.cc.Invoke(ctx, "/proto.UserService/GetWalletBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	RegisterUser(context.Context, *UserRegistrationPayload) (*RegisterMessage, error)
	LoginUser(context.Context, *UserLoginPayload) (*LoginResponse, error)
	CreateRating(context.Context, *ReviewRequest) (*ReviewMessage, error)
	GetUserByID(context.Context, *UserIDMessage) (*User, error)
	ChargeUserWallet(context.Context, *ChargeRequest) (*ChargeResponse, error)
	CreditUserWallet(context.Context, *ChargeRequest) (*ChargeResponse, error)
	VerifyOTP(context.Context, *OTPVerifyPayload) (*OTPVerifyResponse, error)
	ResendOTP(context.Context, *OTPResendPayload) (*OTPResendResponse, error)
	GetWalletBalance(context.Context, *WalletBalanceRequest) (*WalletBalanceResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) RegisterUser(context.Context, *UserRegistrationPayload) (*RegisterMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (UnimplementedUserServiceServer) LoginUser(context.Context, *UserLoginPayload) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedUserServiceServer) CreateRating(context.Context, *ReviewRequest) (*ReviewMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRating not implemented")
}
func (UnimplementedUserServiceServer) GetUserByID(context.Context, *UserIDMessage) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByID not implemented")
}
func (UnimplementedUserServiceServer) ChargeUserWallet(context.Context, *ChargeRequest) (*ChargeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChargeUserWallet not implemented")
}
func (UnimplementedUserServiceServer) CreditUserWallet(context.Context, *ChargeRequest) (*ChargeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreditUserWallet not implemented")
}
func (UnimplementedUserServiceServer) VerifyOTP(context.Context, *OTPVerifyPayload) (*OTPVerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyOTP not implemented")
}
func (UnimplementedUserServiceServer) ResendOTP(context.Context, *OTPResendPayload) (*OTPResendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResendOTP not implemented")
}
func (UnimplementedUserServiceServer) GetWalletBalance(context.Context, *WalletBalanceRequest) (*WalletBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWalletBalance not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRegistrationPayload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/RegisterUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RegisterUser(ctx, req.(*UserRegistrationPayload))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserLoginPayload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/LoginUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).LoginUser(ctx, req.(*UserLoginPayload))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CreateRating_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateRating(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/CreateRating",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateRating(ctx, req.(*ReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIDMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/GetUserByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserByID(ctx, req.(*UserIDMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ChargeUserWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChargeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ChargeUserWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/ChargeUserWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ChargeUserWallet(ctx, req.(*ChargeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CreditUserWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChargeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreditUserWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/CreditUserWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreditUserWallet(ctx, req.(*ChargeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_VerifyOTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OTPVerifyPayload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).VerifyOTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/VerifyOTP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).VerifyOTP(ctx, req.(*OTPVerifyPayload))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ResendOTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OTPResendPayload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ResendOTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/ResendOTP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ResendOTP(ctx, req.(*OTPResendPayload))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetWalletBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WalletBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetWalletBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/GetWalletBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetWalletBalance(ctx, req.(*WalletBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterUser",
			Handler:    _UserService_RegisterUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _UserService_LoginUser_Handler,
		},
		{
			MethodName: "CreateRating",
			Handler:    _UserService_CreateRating_Handler,
		},
		{
			MethodName: "GetUserByID",
			Handler:    _UserService_GetUserByID_Handler,
		},
		{
			MethodName: "ChargeUserWallet",
			Handler:    _UserService_ChargeUserWallet_Handler,
		},
		{
			MethodName: "CreditUserWallet",
			Handler:    _UserService_CreditUserWallet_Handler,
		},
		{
			MethodName: "VerifyOTP",
			Handler:    _UserService_VerifyOTP_Handler,
		},
		{
			MethodName: "ResendOTP",
			Handler:    _UserService_ResendOTP_Handler,
		},
		{
			MethodName: "GetWalletBalance",
			Handler:    _UserService_GetWalletBalance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common/proto/users/users.proto",
}
