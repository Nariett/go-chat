package handler

import (
	"Server/internal/storage/repos/activity"
	"Server/internal/storage/repos/message"
	"Server/internal/storage/repos/user"
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"sync"
)

type Handler interface {
	proto.ChatServiceServer
	RegisterUser(ctx context.Context, user *proto.UserData) (*proto.ServerResponse, error)
	AuthenticateUser(ctx context.Context, user *proto.UserData) (*proto.ServerResponse, error)
	JoinChat(user *proto.UserName, stream proto.ChatService_JoinChatServer) error
	LeaveChat(_ context.Context, user *proto.UserName) (*proto.ServerResponse, error)
	SendMessage(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error)
	MarkMessageAsRead(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error)
	MarkAllMessagesAsRead(_ context.Context, id *proto.UserId) (*proto.ServerResponse, error)
	InsertMessage(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error)
	GetUsersActivity(_ context.Context, _ *proto.Empty) (*proto.Users, error)
	GetUsers(_ context.Context, _ *proto.Empty) (*proto.Users, error)
	GetUserId(_ context.Context, user *proto.UserName) (*proto.UserId, error)
	GetUnreadMessageCount(_ context.Context, userId *proto.UserId) (*proto.UnreadMessages, error)
	GetUsersActivityDates(_ context.Context, empty *proto.Empty) (*proto.UserActivityDates, error)
	GetUnreadMessagesFromUser(_ context.Context, user *proto.UnreadChat) (*proto.UserMessages, error)
	MarkAllMessagesAsReadFromUser(context.Context, *proto.UnreadChat) (*proto.ServerResponse, error)
}

type handler struct {
	proto.UnimplementedChatServiceServer
	mu       sync.Mutex
	users    map[string]chan proto.UserMessage
	user     user.Store
	message  message.Store
	activity activity.Store
}

func NewHandler(user user.Store, message message.Store, activity activity.Store) Handler {
	return &handler{
		users:    make(map[string]chan proto.UserMessage),
		user:     user,
		message:  message,
		activity: activity,
	}
}
