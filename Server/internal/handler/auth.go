package handler

import (
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	"log"
)

func (h *handler) RegUser(ctx context.Context, user *proto.UserData) (*proto.ServerResponse, error) {
	resultChan := make(chan *proto.ServerResponse)
	errorChan := make(chan error)
	go func() {
		response, err := h.user.InsertUser(user)
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

func (h *handler) AuthUser(ctx context.Context, user *proto.UserData) (*proto.ServerResponse, error) {
	resultChan := make(chan *proto.ServerResponse)
	errorChan := make(chan error)
	go func() {
		response, err := h.user.GetUserIdWithUpdateActivity(user)
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
