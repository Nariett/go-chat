package server

import (
	reposActivity "Server/internal/storage/repos/activity"
	reposMessage "Server/internal/storage/repos/message"
	repos "Server/internal/storage/repos/user"
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"log"
	"strings"
)

func (c *ChatServer) GetActiveUsers(_ context.Context, _ *proto.Empty) (*proto.Users, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var activeUsers []string
	for key := range c.users {
		activeUsers = append(activeUsers, key)
	}
	log.Println("Активные пользователи:", strings.Join(activeUsers, " "))
	return &proto.Users{Usernames: activeUsers}, nil
}

func (c *ChatServer) GetUsers(_ context.Context, _ *proto.Empty) (*proto.Users, error) {
	users, err := repos.GetUsers(c.db)
	if err != nil {
		log.Fatal("Ошибка получения данных")
	}
	return &proto.Users{Usernames: users}, nil
}
func (c *ChatServer) GetUserId(_ context.Context, user *proto.UserName) (*proto.UserId, error) {
	userId, err := repos.GetUserId(c.db, user.Name)
	if err != nil {
		log.Fatalf("Ошибка получения id %v", err)
	}
	return &proto.UserId{Id: userId}, err
}
func (c *ChatServer) GetUnreadMessagesCounter(_ context.Context, userId *proto.UserId) (*proto.UnreadMessages, error) {
	unreadMessages, err := reposMessage.GetUnreadMessagesCounter(c.db, userId)
	if err != nil {
		log.Fatal("Ошибка получения непрочитанных сообщений")
	}
	return unreadMessages, nil
}

func (c *ChatServer) GetUsersActivityDates(_ context.Context, empty *proto.Empty) (*proto.UserActivityDates, error) {
	userActivityDates, err := reposActivity.GetUsersActivityDates(c.db, empty)
	if err != nil {
		log.Fatalf("Ошибка получения последней активности пользователей%v", err)
	}
	return userActivityDates, nil
}
