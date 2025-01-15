package handler

import (
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"log"
)

func (h *handler) ReadOneMessage(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error) {
	go func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		err := h.message.UpdateMessageReadStatus(msg)
		if err != nil {
			log.Fatalf("Ошибка чтения сообщения %v", err)
		}
	}()
	return &proto.Empty{}, nil
}

func (h *handler) ReadAllMessages(_ context.Context, id *proto.UserId) (*proto.ServerResponse, error) {
	go func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		err := h.message.UpdateAllMessageReadStatus(id)
		if err != nil {
			log.Fatalf("Ошибка чтения всех сообщения %v", err)
		}
	}()
	return &proto.ServerResponse{Success: true, Message: "Все сообщения прочтены."}, nil
}

func (h *handler) InsertMessage(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error) {
	go func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		err := h.message.InsertMessage(msg)
		if err != nil {
			log.Fatalf("Ошибка добавления сообщения %v", err)
		}
	}()
	return &proto.Empty{}, nil
}
