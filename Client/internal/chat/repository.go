package chat

import (
	"context"

	proto "github.com/Nariett/go-chat/Proto"
)

type ChatRepository struct {
	client proto.ChatServiceClient
}

func NewChatRepository(client proto.ChatServiceClient) *ChatRepository {
	return &ChatRepository{client: client}
}

func (r *ChatRepository) AuthUser(name, password string) (*proto.ServerResponse, error) {
	return r.client.AuthUser(context.Background(), &proto.UserData{Name: name, Password: password})
}

func (r *ChatRepository) RegUser(name, password string) (*proto.ServerResponse, error) {
	return r.client.RegUser(context.Background(), &proto.UserData{Name: name, Password: password})
}
