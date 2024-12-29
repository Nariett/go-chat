package main

import (
	"Client/config"
	"Client/internal/chat"
	"fmt"
	"google.golang.org/grpc"
	"log"

	proto "github.com/Nariett/go-chat/Proto"
)

func main() {

	config := config.LoadConfig()

	connStr := config.BuildConnStr()

	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close()

	client := chat.NewChatRepository(proto.NewChatServiceClient(conn))

	name := chat.InitUser(client)

	stream, err := client.JoinChat(name)
	if err != nil {
		log.Fatalf("Ошибка подключения к чату: %v", err)
	}

	users, err := client.GetUsers(name)
	if err != nil {
		log.Fatalf("Ошибка получения списка пользователй: %v", err)
	}
	fmt.Println("Список всех пользователей:", users.Usernames)

	go client.ListenChat(stream)
	fmt.Println("Для выхода из чата введите \"Выйти\"")
	for {
		var recipient, message string

		fmt.Println("Введите имя, кому хотите отправить сообщение: ")
		fmt.Scanln(&recipient)
		if recipient == "Выйти" {
			chat.ExitChat(client, name)
		}
		fmt.Println("Введите сообщение: ")
		fmt.Scanln(&message)
		if message == "Выйти" {
			chat.ExitChat(client, name)
		}
		if len(recipient) != 0 && len(message) != 0 {
			_, err := client.SendMessage(name, recipient, message)
			if err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
		} else {
			fmt.Println("Сообщение не отправлено. Введите имя пользователя и сообщение.")
		}
	}
}
