package main

import (
	"Client/config"
	"Client/internal/chat"
	"fmt"
	"log"

	"google.golang.org/grpc"

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

	for {
		var recipient, message string

		fmt.Println("Введите имя, кому хотите отправить сообщение: ")
		fmt.Scanln(&recipient)

		fmt.Println("Введите сообщение: ")
		fmt.Scanln(&message)

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
