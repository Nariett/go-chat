package chat

import (
	"bufio"
	"context"
	"fmt"
	proto "github.com/Nariett/go-chat/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"strings"
	"time"
)

func (r *ChatRepository) GetActiveUsers() (*proto.Users, error) {
	return r.client.GetActiveUsers(context.Background(), &proto.Empty{})
}
func (r *ChatRepository) GetUsers(name string) (*proto.Users, error) {
	return r.client.GetUsers(context.Background(), &proto.User{Name: name})
}
func (r *ChatRepository) GetUnreadMessages(name string) (*proto.UnreadMessages, error) {
	return r.client.GetUnreadMessages(context.Background(), &proto.User{Name: name})
}
func (r *ChatRepository) GetUsersActivityDates() (*proto.UserActivityDates, error) {
	return r.client.GetUsersActivityDates(context.Background(), &proto.Empty{})
}
func (r *ChatRepository) JoinChat(name string) (proto.ChatService_JoinChatClient, error) {
	return r.client.JoinChat(context.Background(), &proto.User{Name: name})
}
func (r *ChatRepository) LeaveChat(name string) (*proto.ServerResponse, error) {
	return r.client.LeaveChat(context.Background(), &proto.User{Name: name})
}

func (r *ChatRepository) SendMessage(sender, recipient, content string) (*proto.Empty, error) {
	message := &proto.UserMessage{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
		SentAt:    timestamppb.New(time.Now()),
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
		if msg.Sender == r.CurrentChatUser {
			fmt.Printf("[%s]: %s\n", msg.Sender, msg.Content)
		}
	}
}

func (r *ChatRepository) GetOnlineUsersWithMessageCount(name string) []string {
	activeUsers, err := r.GetActiveUsers()
	if err != nil {
		log.Fatalf("Ошибка получения списка активных пользователй: %v", err)
	}

	users, err := r.GetUsers(name)
	if err != nil {
		log.Fatalf("Ошибка получения списка пользователй: %v", err)
	}

	messageCount, err := r.GetUnreadMessages(name)
	if err != nil {
		log.Fatalf("Ошибка получения списка полученных сообщений: %v", err)
	}

	usersActivityDates, err := r.GetUsersActivityDates()
	if err != nil {
		log.Fatalf("Ошибка получения списка последней активности пользователей: %v", err)
	}
	var allUsers []string

	for _, user := range users.Usernames {
		count := messageCount.Messages[user]
		activityTime := usersActivityDates.ActivityDate[user].AsTime()
		formattedTime := activityTime.Format("15:04:05 02.01.2006") + " - последняя активность"
		status := ""
		if Contains(activeUsers, user) {
			status = " *"
			formattedTime = ""
		}
		if count > 0 {
			allUsers = append(allUsers, fmt.Sprintf("%s (%d)%s", user, count, status))
		} else {
			allUsers = append(allUsers, fmt.Sprintf("%s%s\t%s", user, status, formattedTime))
		}
	}

	return allUsers
}

func InitUser(client *ChatRepository) string {
	var (
		name     string
		password string
		flag     bool = false
	)
	scanner := bufio.NewScanner(os.Stdin) // Создаём сканер для ввода
	for {
		fmt.Println("1 - Войти в чат\n2 - Зарегистрироваться в чате\n3 - Выйти из чата")
		scanner.Scan() // Читаем ввод как строку
		value := scanner.Text()
		switch value {
		case "1":
			fmt.Println("Введите имя: ")
			scanner.Scan()
			name = scanner.Text()
			fmt.Println("Введите пароль: ")
			scanner.Scan()
			password = scanner.Text()
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
		case "2":
			for {
				fmt.Println("Введите имя: ")
				scanner.Scan()
				name := scanner.Text()
				fmt.Println("Введите пароль: ")
				scanner.Scan()
				password = scanner.Text()
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

		case "3":
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

func ExitChat(client *ChatRepository, name string) {
	response, err := client.LeaveChat(name)
	if err != nil {
		fmt.Println("Ошибка)")
	}
	fmt.Println(response.Message)
	os.Exit(0)
}

func Contains(users *proto.Users, username string) bool {
	for _, u := range users.Usernames {
		if u == username {
			return true
		}
	}
	return false
}
func ArrayContainsSubstring(stringArray []string, stringCheck string) bool {
	for _, value := range stringArray {
		if strings.Contains(value, stringCheck) {
			return true
		}
	}
	return false
}
