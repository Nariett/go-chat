package main

import (
	"Server/config"
	"Server/internal/handler"
	"Server/internal/storage"
	"Server/internal/storage/repos/activity"
	"Server/internal/storage/repos/message"
	"Server/internal/storage/repos/user"
	"Server/schema"
	"context"
	proto "github.com/Nariett/go-chat/Proto"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func StartServer(lc fx.Lifecycle, h handler.Handler, conf *config.Config) {
	protocol, port := conf.GetProtocolAndPort()
	listener, err := net.Listen(protocol, port)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	proto.RegisterChatServiceServer(server, h)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("gRPC-сервер запущен")
			go func() {
				if err := server.Serve(listener); err != nil {
					log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.GracefulStop()
			log.Println("gRPC-сервер остановлен")
			return nil
		},
	})
}

func main() {
	application := fx.New(
		fx.Provide(
			config.NewConfig,
			storage.CreatePostgresConnection,
			user.NewStore,
			message.NewStore,
			activity.NewStore,
			storage.Construct,
			handler.NewHandler,
		),
		fx.Invoke(
			schema.Migrate,
			StartServer),
	)
	application.Run()
}
