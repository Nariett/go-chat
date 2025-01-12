package server

import (
	reposActivity "Server/internal/storage/repos/activity"
	repos "Server/internal/storage/repos/user"
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"log"
)

func (c *ChatServer) JoinChat(user *proto.UserName, stream proto.ChatService_JoinChatServer) error {
	c.mu.Lock()
	msgChan := make(chan proto.UserMessage, 10)
	c.users[user.Name] = msgChan
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.users, user.Name)
		close(msgChan)
		c.mu.Unlock()
	}()

	for msg := range msgChan {
		if err := stream.Send(&msg); err != nil {
			log.Printf("Ошибка отправки сообщения клиенту %s: %v", user.Name, err)
			return err
		}
	}
	return nil
}

func (c *ChatServer) LeaveChat(_ context.Context, user *proto.UserName) (*proto.ServerResponse, error) {
	c.mu.Lock()
	userId, err := repos.GetUserId(c.db, user.Name)
	if err != nil {
		log.Fatal(err)
	}
	err = reposActivity.UpdateLastActivity(c.db, userId)
	if err != nil {
		log.Fatal(err)
	}
	delete(c.users, user.Name)
	c.mu.Unlock()
	log.Printf("Пользователь %s вышел из чата", user.Name)
	return &proto.ServerResponse{
		Success: true,
		Message: "Вы вышли из чата",
	}, nil
}

func (c *ChatServer) SendMessage(ctx context.Context, msg *proto.UserMessage) (*proto.Empty, error) {
	go func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		if ch, exists := c.users[msg.Recipient]; exists {
			ch <- *msg
		}

		_, err := c.InsertMessage(ctx, msg)
		if err != nil {
			log.Fatal("Ошибка записи сообщения:", err)
		}
	}()
	return &proto.Empty{}, nil
}
