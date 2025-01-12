package chat

import (
	"context"
	proto "github.com/Nariett/go-chat/Proto"
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
func (r *ChatRepository) AuthUser(name, password string) (*proto.ServerResponse, error) {
	return r.client.AuthUser(context.Background(), &proto.UserData{Name: name, Password: password})
}

func (r *ChatRepository) RegUser(name, password string) (*proto.ServerResponse, error) {
	return r.client.RegUser(context.Background(), &proto.UserData{Name: name, Password: password})
}
