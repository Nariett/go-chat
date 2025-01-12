package server

import (
	reposMessage "Server/internal/storage/repos/message"
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"log"
)

func (c *ChatServer) ReadOneMessage(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error) {
	go func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		err := reposMessage.UpdateMessageReadStatus(c.db, msg)
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
		err := reposMessage.UpdateAllMessageReadStatus(c.db, id)
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
		err := reposMessage.InsertMessage(c.db, msg)
		if err != nil {
			log.Fatalf("Ошибка добавления сообщения %v", err)
		}
	}()
	return &proto.Empty{}, nil
}
