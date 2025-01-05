// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0
// source: chat.proto

package Proto

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
	ChatService_RegUser_FullMethodName               = "/chat.ChatService/RegUser"
	ChatService_AuthUser_FullMethodName              = "/chat.ChatService/AuthUser"
	ChatService_JoinChat_FullMethodName              = "/chat.ChatService/JoinChat"
	ChatService_LeaveChat_FullMethodName             = "/chat.ChatService/LeaveChat"
	ChatService_GetUsers_FullMethodName              = "/chat.ChatService/GetUsers"
	ChatService_GetActiveUsers_FullMethodName        = "/chat.ChatService/GetActiveUsers"
	ChatService_GetUsersActivityDates_FullMethodName = "/chat.ChatService/GetUsersActivityDates"
	ChatService_GetUnreadMessages_FullMethodName     = "/chat.ChatService/GetUnreadMessages"
	ChatService_SendMessage_FullMethodName           = "/chat.ChatService/SendMessage"
	ChatService_ReadMessage_FullMethodName           = "/chat.ChatService/ReadMessage"
)

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	RegUser(ctx context.Context, in *UserData, opts ...grpc.CallOption) (*ServerResponse, error)
	AuthUser(ctx context.Context, in *UserData, opts ...grpc.CallOption) (*ServerResponse, error)
	JoinChat(ctx context.Context, in *User, opts ...grpc.CallOption) (grpc.ServerStreamingClient[UserMessage], error)
	LeaveChat(ctx context.Context, in *User, opts ...grpc.CallOption) (*ServerResponse, error)
	GetUsers(ctx context.Context, in *User, opts ...grpc.CallOption) (*Users, error)
	GetActiveUsers(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Users, error)
	GetUsersActivityDates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*UserActivityDates, error)
	GetUnreadMessages(ctx context.Context, in *User, opts ...grpc.CallOption) (*UnreadMessages, error)
	SendMessage(ctx context.Context, in *UserMessage, opts ...grpc.CallOption) (*Empty, error)
	ReadMessage(ctx context.Context, in *UnreadChat, opts ...grpc.CallOption) (*ServerResponse, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) RegUser(ctx context.Context, in *UserData, opts ...grpc.CallOption) (*ServerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, ChatService_RegUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) AuthUser(ctx context.Context, in *UserData, opts ...grpc.CallOption) (*ServerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, ChatService_AuthUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) JoinChat(ctx context.Context, in *User, opts ...grpc.CallOption) (grpc.ServerStreamingClient[UserMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[0], ChatService_JoinChat_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[User, UserMessage]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ChatService_JoinChatClient = grpc.ServerStreamingClient[UserMessage]

func (c *chatServiceClient) LeaveChat(ctx context.Context, in *User, opts ...grpc.CallOption) (*ServerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, ChatService_LeaveChat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetUsers(ctx context.Context, in *User, opts ...grpc.CallOption) (*Users, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Users)
	err := c.cc.Invoke(ctx, ChatService_GetUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetActiveUsers(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Users, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Users)
	err := c.cc.Invoke(ctx, ChatService_GetActiveUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetUsersActivityDates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*UserActivityDates, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserActivityDates)
	err := c.cc.Invoke(ctx, ChatService_GetUsersActivityDates_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetUnreadMessages(ctx context.Context, in *User, opts ...grpc.CallOption) (*UnreadMessages, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UnreadMessages)
	err := c.cc.Invoke(ctx, ChatService_GetUnreadMessages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) SendMessage(ctx context.Context, in *UserMessage, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ChatService_SendMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) ReadMessage(ctx context.Context, in *UnreadChat, opts ...grpc.CallOption) (*ServerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, ChatService_ReadMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility.
type ChatServiceServer interface {
	RegUser(context.Context, *UserData) (*ServerResponse, error)
	AuthUser(context.Context, *UserData) (*ServerResponse, error)
	JoinChat(*User, grpc.ServerStreamingServer[UserMessage]) error
	LeaveChat(context.Context, *User) (*ServerResponse, error)
	GetUsers(context.Context, *User) (*Users, error)
	GetActiveUsers(context.Context, *Empty) (*Users, error)
	GetUsersActivityDates(context.Context, *Empty) (*UserActivityDates, error)
	GetUnreadMessages(context.Context, *User) (*UnreadMessages, error)
	SendMessage(context.Context, *UserMessage) (*Empty, error)
	ReadMessage(context.Context, *UnreadChat) (*ServerResponse, error)
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChatServiceServer struct{}

func (UnimplementedChatServiceServer) RegUser(context.Context, *UserData) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegUser not implemented")
}
func (UnimplementedChatServiceServer) AuthUser(context.Context, *UserData) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthUser not implemented")
}
func (UnimplementedChatServiceServer) JoinChat(*User, grpc.ServerStreamingServer[UserMessage]) error {
	return status.Errorf(codes.Unimplemented, "method JoinChat not implemented")
}
func (UnimplementedChatServiceServer) LeaveChat(context.Context, *User) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveChat not implemented")
}
func (UnimplementedChatServiceServer) GetUsers(context.Context, *User) (*Users, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}
func (UnimplementedChatServiceServer) GetActiveUsers(context.Context, *Empty) (*Users, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActiveUsers not implemented")
}
func (UnimplementedChatServiceServer) GetUsersActivityDates(context.Context, *Empty) (*UserActivityDates, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsersActivityDates not implemented")
}
func (UnimplementedChatServiceServer) GetUnreadMessages(context.Context, *User) (*UnreadMessages, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUnreadMessages not implemented")
}
func (UnimplementedChatServiceServer) SendMessage(context.Context, *UserMessage) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedChatServiceServer) ReadMessage(context.Context, *UnreadChat) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadMessage not implemented")
}
func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}
func (UnimplementedChatServiceServer) testEmbeddedByValue()                     {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServiceServer will
// result in compilation errors.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	// If the following call pancis, it indicates UnimplementedChatServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_RegUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).RegUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_RegUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).RegUser(ctx, req.(*UserData))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_AuthUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).AuthUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_AuthUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).AuthUser(ctx, req.(*UserData))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_JoinChat_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(User)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServiceServer).JoinChat(m, &grpc.GenericServerStream[User, UserMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ChatService_JoinChatServer = grpc.ServerStreamingServer[UserMessage]

func _ChatService_LeaveChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).LeaveChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_LeaveChat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).LeaveChat(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_GetUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetUsers(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetActiveUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetActiveUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_GetActiveUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetActiveUsers(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetUsersActivityDates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetUsersActivityDates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_GetUsersActivityDates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetUsersActivityDates(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetUnreadMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetUnreadMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_GetUnreadMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetUnreadMessages(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).SendMessage(ctx, req.(*UserMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_ReadMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnreadChat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).ReadMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_ReadMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).ReadMessage(ctx, req.(*UnreadChat))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegUser",
			Handler:    _ChatService_RegUser_Handler,
		},
		{
			MethodName: "AuthUser",
			Handler:    _ChatService_AuthUser_Handler,
		},
		{
			MethodName: "LeaveChat",
			Handler:    _ChatService_LeaveChat_Handler,
		},
		{
			MethodName: "GetUsers",
			Handler:    _ChatService_GetUsers_Handler,
		},
		{
			MethodName: "GetActiveUsers",
			Handler:    _ChatService_GetActiveUsers_Handler,
		},
		{
			MethodName: "GetUsersActivityDates",
			Handler:    _ChatService_GetUsersActivityDates_Handler,
		},
		{
			MethodName: "GetUnreadMessages",
			Handler:    _ChatService_GetUnreadMessages_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _ChatService_SendMessage_Handler,
		},
		{
			MethodName: "ReadMessage",
			Handler:    _ChatService_ReadMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "JoinChat",
			Handler:       _ChatService_JoinChat_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chat.proto",
}
