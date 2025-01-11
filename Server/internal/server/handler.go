package server

import (
	"Server/internal/storage"
	"context"
	"database/sql"
	proto "github.com/Nariett/go-chat/Proto"
	"log"
	"strings"
	"sync"
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
	userId, err := storage.GetUserId(c.db, user.Name)
	if err != nil {
		log.Fatal(err)
	}
	err = storage.UpdateLastActivity(c.db, userId)
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
	users, err := storage.GetUsers(c.db)
	if err != nil {
		log.Fatal("Ошибка получения данных")
	}
	return &proto.Users{Usernames: users}, nil
}
func (c *ChatServer) GetUserId(_ context.Context, user *proto.UserName) (*proto.UserId, error) {
	userId, err := storage.GetUserId(c.db, user.Name)
	if err != nil {
		log.Fatalf("Ошибка получения id %v", err)
	}
	return &proto.UserId{Id: userId}, err
}
func (c *ChatServer) GetUnreadMessagesCounter(_ context.Context, userId *proto.UserId) (*proto.UnreadMessages, error) {
	unreadMessages, err := storage.GetUnreadMessagesCounter(c.db, userId)
	if err != nil {
		log.Fatal("Ошибка получения непрочитанных сообщений")
	}
	return unreadMessages, nil
}

func (c *ChatServer) GetUsersActivityDates(_ context.Context, empty *proto.Empty) (*proto.UserActivityDates, error) {
	userActivityDates, err := storage.GetUsersActivityDates(c.db, empty)
	if err != nil {
		log.Fatal("Ошибка получения последней активности пользователей")
	}
	return userActivityDates, nil
}

func (c *ChatServer) ReadOneMessage(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error) {
	go func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		err := storage.ReadOneMessage(c.db, msg)
		if err != nil {
			log.Fatalf("Ошибка чтения сообщения %v", err)
		}
	}()
	return &proto.Empty{}, nil
}
func (c *ChatServer) ReadAllMessages(_ context.Context, id *proto.UserId) (*proto.ServerResponse, error) {
	go func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		err := storage.RealAllMessages(c.db, id)
		if err != nil {
			log.Fatalf("Ошибка чтения всех сообщения %v", err)
		}
	}()
	return &proto.ServerResponse{Success: true, Message: "Все сообщения прочтены."}, nil
}
func (c *ChatServer) InsertMessage(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error) {
	go func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		err := storage.InsertMessage(c.db, msg)
		if err != nil {
			log.Fatalf("Ошибка добавления сообщения %v", err)
		}
	}()
	return &proto.Empty{}, nil
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

func (c *ChatServer) RegUser(ctx context.Context, user *proto.UserData) (*proto.ServerResponse, error) {
	resultChan := make(chan *proto.ServerResponse)
	errorChan := make(chan error)
	go func() {
		response, err := storage.RegUser(c.db, user)
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
		response, err := storage.AuthUser(c.db, user)
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
