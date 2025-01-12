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
)

func (r *ChatRepository) GetActiveUsers() (*proto.Users, error) {
	return r.client.GetActiveUsers(context.Background(), &proto.Empty{})
}
func (r *ChatRepository) GetUsers() (*proto.Users, error) {
	return r.client.GetUsers(context.Background(), &proto.Empty{})
}
func (r *ChatRepository) GetUserId(name string) (int32, error) {
	idUser, err := r.client.GetUserId(context.Background(), &proto.UserName{Name: name})
	if err != nil {
		return -1, err
	}
	return idUser.Id, nil
}
func (r *ChatRepository) GetUnreadMessagesCounter(id int32) (*proto.UnreadMessages, error) {
	return r.client.GetUnreadMessagesCounter(context.Background(), &proto.UserId{Id: id})
}
func (r *ChatRepository) GetUsersActivityDates() (*proto.UserActivityDates, error) {
	return r.client.GetUsersActivityDates(context.Background(), &proto.Empty{})
}
func (r *ChatRepository) JoinChat(name string) (proto.ChatService_JoinChatClient, error) {
	return r.client.JoinChat(context.Background(), &proto.UserName{Name: name})
}
func (r *ChatRepository) LeaveChat(name string) (*proto.ServerResponse, error) {
	return r.client.LeaveChat(context.Background(), &proto.UserName{Name: name})
}
func (r *ChatRepository) ReadAllMessages(id int32) (*proto.ServerResponse, error) {
	return r.client.ReadAllMessages(context.Background(), &proto.UserId{Id: id})
}

func (r *ChatRepository) SendMessage(sender string, senderId int32, recipient string, recipientId int32, content string) (*proto.Empty, error) {
	message := &proto.UserMessage{
		Sender:      sender,
		SenderId:    senderId,
		Recipient:   recipient,
		RecipientId: recipientId,
		Content:     content,
		SentAt:      timestamppb.Now(),
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
			_, err = r.client.ReadOneMessage(context.Background(), msg)
			if err != nil {
				log.Fatalf("Ошибка обновления данных: %v", err)
			}
		}
	}
}

func (r *ChatRepository) GetOnlineUsersWithMessageCount(id int32, name string) []string {
	activeUsers, err := r.GetActiveUsers()
	if err != nil {
		log.Fatalf("Ошибка получения списка активных пользователй: %v", err)
	}

	users, err := r.GetUsers()
	if err != nil {
		log.Fatalf("Ошибка получения списка пользователй: %v", err)
	}

	messageCount, err := r.GetUnreadMessagesCounter(id)
	if err != nil {
		log.Fatalf("Ошибка получения списка полученных сообщений: %v", err)
	}

	usersActivityDates, err := r.GetUsersActivityDates()
	if err != nil {
		log.Fatalf("Ошибка получения списка последней активности пользователей: %v", err)
	}
	var allUsers []string

	for _, user := range users.Usernames {
		if user != name {
			count := messageCount.Messages[user]
			activityTime := usersActivityDates.ActivityDate[user].AsTime().In(r.Location)
			status := activityTime.Format("15:04:05 02.01.2006") + " - последняя активность"
			if Contains(activeUsers, user) {
				status = "В сети"
			}
			if count > 0 {
				allUsers = append(allUsers, fmt.Sprintf("%s (%d)\t%s", user, count, status))
			} else {
				allUsers = append(allUsers, fmt.Sprintf("%s\t\t%s", user, status))
			}
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
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("1 - Войти в чат\n2 - Зарегистрироваться в чате\n3 - Выйти из чата")
		scanner.Scan()
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
		name := strings.Split(value, "\t")[0]
		if strings.Contains(name, stringCheck) {
			return true
		}
	}
	return false
}

func ChatSession(client *ChatRepository, name string, userId int32, recipient string, recipientId int32) error {
	fmt.Println("Открыт чат с пользователем :", recipient, "для выхода в чаты напишите '/Чаты'")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		message := scanner.Text()
		if message == "/Чаты" {
			fmt.Println("Вы перешли в чаты")
			client.CurrentChatUser = ""
			break
		}
		if len(recipient) != 0 && len(message) != 0 {
			_, err := client.SendMessage(name, userId, recipient, recipientId, message)
			if err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
		} else {
			fmt.Println("Сообщение не отправлено. Введите имя пользователя и сообщение.")
		}
	}
	return nil
}
