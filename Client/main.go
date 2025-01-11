package main

import (
	"Client/config"
	"Client/internal/chat"
	"bufio"
	"fmt"
	proto "github.com/Nariett/go-chat/Proto"
	"google.golang.org/grpc"
	"log"
	"os"
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

	idUser, err := client.GetUserId(name)
	if err != nil {
		log.Fatalf("Ошибка получения id пользователя: %v", err)
	}
	go client.ListenChat(stream)
	for {
		fmt.Println("Для выхода из приложения введите '/Выход'")
		fmt.Println("Для выхода из приложения введите '/Прочитать все'")
		fmt.Println("Список пользователей '(x)' - число новых сообщений")

		onlineUser := client.GetOnlineUsersWithMessageCount(idUser, name)
		for _, user := range onlineUser {
			fmt.Println(user)
		}

		restart := false
		var recipient string
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Println("Выберите чат:")
			scanner.Scan()
			recipient = scanner.Text()
			if recipient != "" && chat.ArrayContainsSubstring(onlineUser, recipient) {
				break
			} else if recipient == "/Выход" {
				chat.ExitChat(client, name)
				fmt.Println("Вы вышли из чата")
				os.Exit(0)
			} else if recipient == "/Прочитать все" {
				response, err := client.ReadAllMessages(idUser)
				if err != nil {
					log.Fatalln("Ошибка чтения сообщений")
				}
				fmt.Println(response.Message)
				restart = true
				break
			} else {
				fmt.Println("Введите корректное имя и повторите попытку")
				continue
			}
		}
		if restart {
			continue
		}
		recipientId, err := client.GetUserId(recipient)
		if err != nil {
			log.Fatalf("Ошибка получения id получателя: %v", err)
		}
		client.CurrentChatUser = recipient
		err = chat.ChatSession(client, name, idUser, recipient, recipientId)
		if err != nil {
			log.Fatalf("Ошибка чата с пользователем: %v", err)
		}
	}
}
