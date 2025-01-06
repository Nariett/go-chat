package storage

import (
	"database/sql"
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"

	proto "github.com/Nariett/go-chat/Proto"
	"github.com/lib/pq"
)

func RegUser(db *sql.DB, user *proto.UserData) (*proto.ServerResponse, error) {
	log.Printf("Добавление пользователя %s,%s", user.Name, user.Password)

	tx, err := db.Begin()
	if err != nil {
		return &proto.ServerResponse{
			Success: false,
			Message: "Ошибка начала транзакции",
		}, err
	}

	var userId int64

	err = tx.QueryRow("INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id", user.Name, user.Password).Scan(&userId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Printf("Данный ник уже занят: %v\n", err)
			_ = tx.Rollback()
			return &proto.ServerResponse{
				Success: false,
				Message: "Пользователь не добавлен в базу данных, так как ник уже занят",
			}, nil
		}
		log.Println("Ошибка транзакции")
		_ = tx.Rollback()
		return &proto.ServerResponse{
			Success: false,
			Message: "Пользователь не добавлен в базу данных, ошибка базы данных, повторите попытку",
		}, nil
	}
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("INSERT INTO activity (idUser,date) VALUES ((SELECT id FROM users WHERE name = $1),$2)", user.Name, time.Now().In(location))
	if err != nil {
		_ = tx.Rollback()
		return &proto.ServerResponse{
			Success: false,
			Message: "Пользователь не добавлен в базу данных, ошибка базы данных, повторите попытку",
		}, nil
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Ошибка транзакции")
		return &proto.ServerResponse{
			Success: false,
			Message: "Ошибка при сохранении изменений в базе данных",
		}, err
	}

	log.Printf("Добавлен новый пользователь: id: %d, name: %s, password: %s\n", userId, user.Name, user.Password)
	return &proto.ServerResponse{
		Success: true,
		Message: "Пользователь добавлен в базу данных",
	}, nil
}
func AuthUser(db *sql.DB, user *proto.UserData) (*proto.ServerResponse, error) {
	log.Printf("Найти пользователя %s, %s", user.Name, user.Password)

	var idUser int
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Ошибка начала транзакции %v\n", err)
	}

	err = tx.QueryRow("SELECT id FROM users WHERE name = $1 AND password = $2", user.Name, user.Password).Scan(&idUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Пользователь не найден")
			_ = tx.Rollback()
			return &proto.ServerResponse{
				Success: false,
				Message: "Пользователь не найден. Проверьте данные и повторите попытку",
			}, nil
		} else {
			_ = tx.Rollback()
			log.Printf("Ошибка выполнения запроса: %v", err)
			return &proto.ServerResponse{
				Success: false,
				Message: "Ошибка получения данных, повторите попытку.",
			}, nil
		}
	}
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("UPDATE activity SET date = $1 WHERE idUser = $2", time.Now().In(location), idUser)
	if err != nil {
		_ = tx.Rollback()
		log.Printf("Ошибка выполнения запроса: %v", err)
		return &proto.ServerResponse{
			Success: false,
			Message: "Ошибка добавления данных, повторите попытку.",
		}, nil
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Ошибка завершения транзакции: %v", err)
	}
	log.Printf("Транзакция выполнена")

	return &proto.ServerResponse{Success: true}, nil
}
func InsertMessage(db *sql.DB, message *proto.UserMessage) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Ошибка начала транзакции %v\n", err)
		return err
	}
	var senderId, recipientId int
	err = tx.QueryRow("SELECT id FROM users WHERE name = $1", message.Sender).Scan(&senderId)
	if err != nil {
		_ = tx.Rollback()
		log.Printf("Ошибка выполнения запроса: %v", err)
		return err
	}
	err = tx.QueryRow("SELECT id FROM users WHERE name = $1", message.Recipient).Scan(&recipientId)
	if err != nil {
		_ = tx.Rollback()
		log.Printf("Ошибка выполнения запроса: %v", err)
		return err
	}

	_, err = tx.Exec("INSERT INTO messages (sender_id, recipient_id, content, sent_at) VALUES ($1, $2, $3, $4)", senderId, recipientId, message.Content, message.SentAt.AsTime())
	if err != nil {
		_ = tx.Rollback()
		log.Printf("Ошибка выполнения запроса: %v", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Ошибка выполнения транзакции: %v", err)
		return err
	}
	return nil
}

func UpdateLastActivity(db *sql.DB, user *proto.User) error {
	location, err := time.LoadLocation("Europe/Moscow")
	_, err = db.Exec("UPDATE activity SET date = $1 WHERE idUser = (SELECT id FROM users WHERE name = $2)", time.Now().In(location), user.Name)
	if err != nil {
		log.Fatalf("Ошибка обновления данных %v\n", err)
	}
	return nil
}

func GetUsers(db *sql.DB, user *proto.User) ([]string, error) {
	rows, err := db.Query("SELECT name FROM users WHERE name != $1", user.Name)
	if err != nil {
		log.Fatalf("Ошибка получения данных %v\n", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal("Ошибка")
		}
	}(rows)

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			log.Fatal("Ошибка чтения строки")
		}
		usernames = append(usernames, username)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Ошибка перебора строк")
	}
	return usernames, nil
}

func GetUnreadMessages(db *sql.DB, user *proto.User) (*proto.UnreadMessages, error) {
	var receiverId int
	err := db.QueryRow("SELECT id FROM users WHERE name = $1", user.Name).Scan(&receiverId)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return &proto.UnreadMessages{Messages: nil}, err
	}
	rows, err := db.Query("SELECT u.name, COUNT(m.id) FROM users u LEFT JOIN messages m ON u.id = m.sender_id AND m.recipient_id = $1 GROUP BY u.id", receiverId)
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		if err != nil {
			log.Fatal("Ошибка")
		}
	}(rows)

	UnreadMessages := &proto.UnreadMessages{
		Messages: make(map[string]int32),
	}

	for rows.Next() {
		var username string
		var count int
		if err := rows.Scan(&username, &count); err != nil {
			log.Fatal("Ошибка чтения строки")
		}
		UnreadMessages.Messages[username] = int32(count)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Ошибка обработки строки")
	}
	return UnreadMessages, nil
}

func GetUsersActivityDates(db *sql.DB, _ *proto.Empty) (*proto.UserActivityDates, error) {
	rows, err := db.Query("SELECT users.name, date FROM public.activity JOIN users on users.id = activity.idUser")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		if err != nil {
			log.Fatal("Ошибка")
		}
	}(rows)
	userActivityDates := &proto.UserActivityDates{
		ActivityDate: make(map[string]*timestamppb.Timestamp),
	}
	for rows.Next() {
		var username string
		var readAt time.Time
		if err := rows.Scan(&username, &readAt); err != nil {
			log.Fatal("Ошибка чтения строки")
		}
		timestamp := timestamppb.New(readAt)
		userActivityDates.ActivityDate[username] = timestamp
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Ошибка обработки строки")
	}
	return userActivityDates, nil
}
