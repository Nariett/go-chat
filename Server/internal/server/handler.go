package server

import (
	"Server/internal/database"
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"

	proto "github.com/Nariett/go-chat/Proto"
)

type ChatServer struct {
	proto.UnimplementedChatServiceServer
	mu    sync.Mutex
	users map[string]chan proto.UserMessage
	db    *sql.DB
}

func newChatServer(db *sql.DB) *ChatServer {
	return &ChatServer{
		users: make(map[string]chan proto.UserMessage),
		db:    db,
	}
}
func (c *ChatServer) JoinChat(user *proto.User, stream proto.ChatService_JoinChatServer) error {
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

func (c *ChatServer) LeaveChat(ctx context.Context, user *proto.User) (*proto.ServerResponse, error) {
	c.mu.Lock()
	err := database.UpdateLastActivity(c.db, user)
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

func (c *ChatServer) GetUsers(ctx context.Context, user *proto.User) (*proto.ActiveUsers, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var activeUsers []string
	for key := range c.users {
		activeUsers = append(activeUsers, key)
	}
	log.Println("Активные пользователи:", strings.Join(activeUsers, " "))
	return &proto.ActiveUsers{Usernames: activeUsers}, nil
}

func (c *ChatServer) SendMessage(ctx context.Context, msg *proto.UserMessage) (*proto.Empty, error) {
	go func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		if ch, exists := c.users[msg.Recipient]; exists {
			ch <- *msg
		}
	}()
	return &proto.Empty{}, nil
}

func (c *ChatServer) RegUser(ctx context.Context, user *proto.UserData) (*proto.ServerResponse, error) {
	resultChan := make(chan *proto.ServerResponse)
	errorChan := make(chan error)
	go func() {
		response, err := database.RegUser(c.db, user)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- response
	}()
	select {
	case response := <-resultChan:
		return response, nil
	case err := <-errorChan:
		log.Printf("Ошибка при регистрации: %v", err)
		return nil, err
	case <-ctx.Done():
		log.Printf("Контекст завершен: %v", ctx.Err())
		return nil, ctx.Err()
	}
}

func (c *ChatServer) AuthUser(ctx context.Context, user *proto.UserData) (*proto.ServerResponse, error) {
	resultChan := make(chan *proto.ServerResponse)
	errorChan := make(chan error)
	go func() {
		response, err := database.AuthUser(c.db, user)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- response
	}()
	select {
	case response := <-resultChan:
		return response, nil
	case err := <-errorChan:
		log.Printf("Ошибка авторизации: %v", err)
		return nil, err
	case <-ctx.Done():
		log.Printf("Контекст завершен: %v", ctx.Err())
		return nil, ctx.Err()
	}
}
