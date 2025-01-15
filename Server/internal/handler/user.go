package handler

import (
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"log"
	"strings"
)

func (h *handler) GetActiveUsers(_ context.Context, _ *proto.Empty) (*proto.Users, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	var activeUsers []string
	for key := range h.users {
		activeUsers = append(activeUsers, key)
	}
	log.Println("Активные пользователи:", strings.Join(activeUsers, " "))
	return &proto.Users{Usernames: activeUsers}, nil
}

func (h *handler) GetUsers(_ context.Context, _ *proto.Empty) (*proto.Users, error) {
	users, err := h.user.GetUsers()
	if err != nil {
		log.Fatal("Ошибка получения данных")
	}
	return &proto.Users{Usernames: users}, nil
}
func (h *handler) GetUserId(_ context.Context, user *proto.UserName) (*proto.UserId, error) {
	userId, err := h.user.GetUserId(user.Name)
	if err != nil {
		log.Fatalf("Ошибка получения id %v", err)
	}
	return &proto.UserId{Id: userId}, err
}
func (h *handler) GetUnreadMessagesCounter(_ context.Context, userId *proto.UserId) (*proto.UnreadMessages, error) {
	unreadMessages, err := h.message.GetUnreadMessagesCounter(userId)
	if err != nil {
		log.Fatal("Ошибка получения непрочитанных сообщений")
	}
	return unreadMessages, nil
}

func (h *handler) GetUsersActivityDates(_ context.Context, empty *proto.Empty) (*proto.UserActivityDates, error) {
	userActivityDates, err := h.activity.GetUsersActivityDates(empty)
	if err != nil {
		log.Fatalf("Ошибка получения последней активности пользователей%v", err)
	}
	return userActivityDates, nil
}
