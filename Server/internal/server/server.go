package server

import (
	"database/sql"
	"log"
	"net"

	proto "github.com/Nariett/go-chat/Proto"

	"google.golang.org/grpc"
)

func StartServer(listener net.Listener, db *sql.DB) {
	server := grpc.NewServer()
	proto.RegisterChatServiceServer(server, newChatServer(db))

	log.Println("gRPC-сервер запущен")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}
}
