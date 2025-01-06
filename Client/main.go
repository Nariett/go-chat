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
	go client.ListenChat(stream)
	for {
		fmt.Println("Для выхода из приложения введите '/Выход'")
		fmt.Println("Список пользователей '*' - в сети, '(x)' - число новых сообщений")
		onlineUser := client.GetOnlineUsersWithMessageCount(name)
		for _, user := range onlineUser {
			fmt.Println(user)
		}
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
			} else {
				fmt.Println("Введите корректное имя и повторите попытку")
				continue
			}
		}
		client.CurrentChatUser = recipient
		fmt.Println("Открыт чат с пользователем :", recipient, "для выхода в чаты напишите '/Чаты'")
		for {
			scanner.Scan()
			message := scanner.Text()
			if message == "/Чаты" {
				fmt.Println("Вы перешли в чаты")
				client.CurrentChatUser = ""
				break
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
