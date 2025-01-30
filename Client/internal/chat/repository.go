package chat

import (
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ChatRepository struct {
	IdUser          int32
	client          proto.ChatServiceClient
	CurrentChatUser string
	Location        *time.Location
}

func NewChatRepository(client proto.ChatServiceClient) *ChatRepository {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic("Не удалось загрузить временную зону Москвы: " + err.Error())
	}
	return &ChatRepository{
		client:          client,
		CurrentChatUser: "",
		Location:        location,
	}
}
func (r *ChatRepository) RegisterUser(name, password string) (*proto.ServerResponse, error) {
	return r.client.RegisterUser(context.Background(), &proto.UserData{Name: name, Password: password})
}

func (r *ChatRepository) AuthenticateUser(name, password string) (*proto.ServerResponse, error) {
	return r.client.AuthenticateUser(context.Background(), &proto.UserData{Name: name, Password: password})
}

func (r *ChatRepository) JoinChat(name string) (proto.ChatService_JoinChatClient, error) {
	return r.client.JoinChat(context.Background(), &proto.UserName{Name: name})
}

func (r *ChatRepository) LeaveChat(name string) (*proto.ServerResponse, error) {
	return r.client.LeaveChat(context.Background(), &proto.UserName{Name: name})
}

func (r *ChatRepository) SendMessage(sender string, senderId int32, recipient string, recipientId int32, content string) (*proto.Empty, error) {
	message := &proto.UserMessage{
		Sender:      sender,
		SenderId:    senderId,
		Recipient:   recipient,
		RecipientId: recipientId,
		Content:     content,
		SentAt:      timestamppb.Now(),
	}
	response, err := r.client.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *ChatRepository) GetUsers() (*proto.Users, error) {
	return r.client.GetUsers(context.Background(), &proto.Empty{})
}

func (r *ChatRepository) GetUsersActivity() (*proto.Users, error) {
	return r.client.GetUsersActivity(context.Background(), &proto.Empty{})
}

func (r *ChatRepository) GetUsersActivityDates() (*proto.UserActivityDates, error) {
	return r.client.GetUsersActivityDates(context.Background(), &proto.Empty{})
}

func (r *ChatRepository) GetUserId(name string) (int32, error) {
	idUser, err := r.client.GetUserId(context.Background(), &proto.UserName{Name: name})
	if err != nil {
		return -1, err
	}
	return idUser.Id, nil
}

func (r *ChatRepository) GetUnreadMessageCount(id int32) (*proto.UnreadMessages, error) {
	return r.client.GetUnreadMessageCount(context.Background(), &proto.UserId{Id: id})
}

func (r *ChatRepository) GetUnreadMessagesFromUser(senderId, recipientId int32) (*proto.UserMessages, error) {
	return r.client.GetUnreadMessagesFromUser(context.Background(), &proto.UnreadChat{Sender: senderId, Recipient: recipientId})
}

func (r *ChatRepository) MarkMessagesAsRead(msg *proto.UserMessage) (*proto.Empty, error) {
	return r.client.MarkMessageAsRead(context.Background(), msg)
}

func (r *ChatRepository) MarkAllMessagesAsRead(id int32) (*proto.ServerResponse, error) {
	return r.client.MarkAllMessagesAsRead(context.Background(), &proto.UserId{Id: id})
}
func (r *ChatRepository) MarkAllMessagesAsReadFromUser(senderId, recipientId int32) (*proto.ServerResponse, error) {
	return r.client.MarkAllMessagesAsReadFromUser(context.Background(), &proto.UnreadChat{Sender: senderId, Recipient: recipientId})
}
