package handler

import (
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"log"
)

func (h *handler) MarkMessageAsRead(_ context.Context, msg *proto.UserMessage) (*proto.Empty, error) {
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

func (h *handler) MarkAllMessagesAsRead(_ context.Context, id *proto.UserId) (*proto.ServerResponse, error) {
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
func (h *handler) MarkAllMessagesAsReadFromUser(_ context.Context, unreadChat *proto.UnreadChat) (*proto.ServerResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	err := h.message.UpdateAllMessagesFromUserReadStatus(unreadChat)
	if err != nil {
		log.Printf("Ошибка получения данных из бд: %v", err)
		return nil, err
	}
	return &proto.ServerResponse{Success: true, Message: "Сообщения от пользователя прочитаны"}, nil
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

func (h *handler) GetUnreadMessagesFromUser(_ context.Context, user *proto.UnreadChat) (*proto.UserMessages, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	result, err := h.message.GetUnreadMessagesFromUser(user)
	if err != nil {
		log.Printf("Ошибка получения данных из бд: %v", err)
		return nil, err
	}
	return result, nil
}
