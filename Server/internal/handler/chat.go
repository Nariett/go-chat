package handler

import (
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"log"
)

func (h *handler) JoinChat(user *proto.UserName, stream proto.ChatService_JoinChatServer) error {
	h.mu.Lock()
	msgChan := make(chan proto.UserMessage, 10)
	h.users[user.Name] = msgChan
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.users, user.Name)
		close(msgChan)
		h.mu.Unlock()
	}()

	for msg := range msgChan {
		if err := stream.Send(&msg); err != nil {
			log.Printf("Ошибка отправки сообщения клиенту %s: %v", user.Name, err)
			return err
		}
	}
	return nil
}

func (h *handler) LeaveChat(_ context.Context, user *proto.UserName) (*proto.ServerResponse, error) {
	h.mu.Lock()
	userId, err := h.user.GetUserId(user.Name)
	if err != nil {
		log.Fatal(err)
	}
	err = h.activity.UpdateLastActivity(userId)
	if err != nil {
		log.Fatal(err)
	}
	delete(h.users, user.Name)
	h.mu.Unlock()
	log.Printf("Пользователь %s вышел из чата", user.Name)
	return &proto.ServerResponse{
		Success: true,
		Message: "Вы вышли из чата",
	}, nil
}

func (h *handler) SendMessage(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error) {
	go func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		if ch, exists := h.users[msg.Recipient]; exists {
			ch <- *msg
		}

		err := h.message.InsertMessage(msg)
		if err != nil {
			log.Fatal("Ошибка записи сообщения:", err)
		}
	}()
	return &proto.Empty{}, nil
}
