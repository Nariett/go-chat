package main

import (
	"Server/config"
	"Server/internal/server"
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
)

func main() {
	loadConfig := config.LoadConfig()

	connStr := loadConfig.BuildConnStr()

	log.Printf("Строка подключения к бд: %s", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	log.Println("База данных подключена")
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Ошибка закрытия подключения: %v", err)
		}
	}()

	protocol, port := loadConfig.GetProtocolAndPort()
	listener, err := net.Listen(protocol, port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	log.Printf("Сервер запущен на порту %s", port)
	server.StartServer(listener, db)
}
