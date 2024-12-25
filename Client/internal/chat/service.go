package chat

import (
	"context"
	"fmt"
	"log"
	"os"

	proto "github.com/Nariett/go-chat/Proto"
)

func (r *ChatRepository) GetUsers(name string) (*proto.ActiveUsers, error) {
	return r.client.GetUsers(context.Background(), &proto.User{Name: name})
}

func (r *ChatRepository) JoinChat(name string) (proto.ChatService_JoinChatClient, error) {
	return r.client.JoinChat(context.Background(), &proto.User{Name: name})
}

func (r *ChatRepository) SendMessage(sender, recipient, content string) (*proto.Empty, error) {
	message := &proto.UserMessage{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
	}
	response, err := r.client.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *ChatRepository) ListenChat(stream proto.ChatService_JoinChatClient) {
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Ошибка получения сообщения: %v", err)
		}
		fmt.Printf("Новое сообщение от %s: %s\n", msg.Sender, msg.Content)
	}
}

func InitUser(client *ChatRepository) string {
	var (
		name     string
		password string
		flag     bool = false
		value    int
	)
	for {
		fmt.Println("1 - Войти в чат\n2 - Зарегистрироваться в чате\n3 - Выйти из чата")
		fmt.Scanln(&value)
		switch value {
		case 1:
			fmt.Println("Введите имя: ")
			fmt.Scanln(&name)
			fmt.Println("Введите пароль: ")
			fmt.Scanln(&password)

			response, err := client.AuthUser(name, password)
			if err != nil {
				log.Fatalf("Ошибка аутентификации: %v", err)
			}
			if response.Success {
				fmt.Println("Вы вошли в систему!")
				flag = true
			} else {
				fmt.Println(response.Message)
			}
		case 2:
			for {
				fmt.Println("Введите имя: ")
				fmt.Scanln(&name)
				fmt.Println("Введите пароль: ")
				fmt.Scanln(&password)
				response, err := client.RegUser(name, password)
				if err != nil {
					log.Fatalf("Ошибка регистрации: %v", err)
				}
				if response.Success {
					fmt.Println("Вы прошли регистрацию!")
					break
				} else {
					fmt.Println(response.Message)
				}
			}

		case 3:
			fmt.Println("Вы вышли из чата...")
			os.Exit(1)
		default:
			fmt.Println("Введите значение и повторите попытку.")
		}
		if flag {
			return name
		}
	}
}
