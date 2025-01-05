package main

import (
	"Client/config"
	"Client/internal/chat"
	"bufio"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"

	proto "github.com/Nariett/go-chat/Proto"
)

func main() {

	loadConfig := config.LoadConfig()

	connStr := loadConfig.BuildConnStr()

	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Ошибка закрытия подключения: %v", err)
		}
	}()

	client := chat.NewChatRepository(proto.NewChatServiceClient(conn))

	name := chat.InitUser(client)

	stream, err := client.JoinChat(name)
	if err != nil {
		log.Fatalf("Ошибка подключения к чату: %v", err)
	}

	go client.ListenChat(stream)
	for {
		onlineUser := client.GetOnlineUsersWithMessageCount(name)
		fmt.Println("Список пользователей (*) - в сети")
		for _, user := range onlineUser {
			fmt.Println(user)
		}
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Выберите чат :")
		scanner.Scan()
		recipientName := scanner.Text()
		fmt.Println("Открыт чат с :", recipientName)
		fmt.Println("Для выхода из чата введите \"Выйти\"")

		//Вывод сообщений из чата

		for {
			var recipient, message string

			fmt.Println("Введите имя, кому хотите отправить сообщение: ")
			scanner.Scan()
			recipient = scanner.Text()
			if recipient == "Выйти" {
				chat.ExitChat(client, name)
			}
			fmt.Println("Введите сообщение: ")
			scanner.Scan()
			message = scanner.Text()
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
}
