package server

import (
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type ChatServer struct {
	proto.UnimplementedChatServiceServer
	mu    sync.Mutex
	users map[string]chan proto.UserMessage
	db    *sqlx.DB
}

func newChatServer(db *sqlx.DB) *ChatServer {
	return &ChatServer{
		users: make(map[string]chan proto.UserMessage),
		db:    db,
	}
}
func StartServer(listener net.Listener, db *sqlx.DB) {
	server := grpc.NewServer()
	proto.RegisterChatServiceServer(server, newChatServer(db))

	log.Println("gRPC-сервер запущен")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}
}
