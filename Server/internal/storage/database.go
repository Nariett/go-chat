package storage

import (
	"Server/internal/storage/dbo"
	"database/sql"
	"errors"
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func RegUser(db *sqlx.DB, user *proto.UserData) (*proto.ServerResponse, error) {
	log.Printf("Добавление пользователя %s,%s", user.Name, user.Password)

	tx, err := db.Beginx()
	if err != nil {
		return &proto.ServerResponse{
			Success: false,
			Message: "Ошибка начала транзакции",
		}, err
	}

	var userId int32

	query := `INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id`
	err = tx.Get(&userId, query, user.Name, user.Password)
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

	query = `INSERT INTO activity (idUser,date) VALUES ((SELECT id FROM users WHERE name = $1),$2)`
	_, err = tx.Exec(query, user.Name, timestamppb.Now().AsTime())
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

func AuthUser(db *sqlx.DB, user *proto.UserData) (*proto.ServerResponse, error) {
	log.Printf("Найти пользователя %s, %s", user.Name, user.Password)

	tx, err := db.Beginx()
	if err != nil {
		log.Fatalf("Ошибка начала транзакции %v\n", err)
	}

	var idUser int

	query := `SELECT id FROM users WHERE name = $1 AND password = $2`
	err = tx.Get(&idUser, query, user.Name, user.Password)
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

	query = `UPDATE activity SET date = $1 WHERE idUser = $2`
	_, err = tx.Exec(query, timestamppb.Now().AsTime(), idUser)
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

func InsertMessage(db *sqlx.DB, message *proto.UserMessage) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Printf("Ошибка начала транзакции %v\n", err)
		return err
	}

	query := `INSERT INTO messages (sender_id, recipient_id, content, sent_at) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(query, message.SenderId, message.RecipientId, message.Content, message.SentAt.AsTime())
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

func UpdateLastActivity(db *sqlx.DB, id int32) error {
	query := `UPDATE activity SET date = $1 WHERE idUser = $2`
	_, err := db.Exec(query, timestamppb.Now().AsTime(), id)
	if err != nil {
		log.Fatalf("Ошибка обновления данных %v\n", err)
	}

	return nil
}

func GetUsers(db *sqlx.DB) ([]string, error) {
	var usernames []string
	err := db.Select(&usernames, "SELECT name FROM users")
	if err != nil {
		return nil, err
	}

	return usernames, nil
}

func GetUserId(db *sqlx.DB, name string) (int32, error) {
	var userId int32
	err := db.Get(&userId, "SELECT id FROM users WHERE name = $1", name)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func GetUnreadMessagesCounter(db *sqlx.DB, id *proto.UserId) (*proto.UnreadMessages, error) {
	var results []dbo.UserUnreadMessageCount
	query := `SELECT u.name, COUNT(m.id) 
			  FROM users u 
			  LEFT JOIN messages m ON u.id = m.sender_id 
			  AND m.recipient_id = $1 
			  AND m.read_at IS NULL GROUP BY u.id`
	err := db.Select(&results, query, id.Id)
	if err != nil {
		return nil, err
	}

	counter := make(map[string]int32)
	for _, result := range results {
		counter[result.Name] = result.Count
	}
	UnreadMessages := &proto.UnreadMessages{
		Messages: counter,
	}

	return UnreadMessages, nil
}

func GetUsersActivityDates(db *sqlx.DB, _ *proto.Empty) (*proto.UserActivityDates, error) {
	var activityDates []dbo.UserActivity
	query := `SELECT users.name AS name, date AS date
			  FROM public.activity 
    		  JOIN users on users.id = activity.idUser`
	err := db.Select(&activityDates, query)
	if err != nil {
		return nil, err
	}

	activityMap := make(map[string]*timestamppb.Timestamp)
	for _, item := range activityDates {
		activityMap[item.Name] = timestamppb.New(item.Date)
	}

	userActivityDates := &proto.UserActivityDates{
		ActivityDate: activityMap,
	}

	return userActivityDates, nil
}

func ReadOneMessage(db *sqlx.DB, message *proto.UserMessage) error {
	query := `UPDATE messages SET read_at = $1 WHERE content = $2 AND sent_at = $3`
	_, err := db.Exec(query, timestamppb.Now().AsTime(), message.Content, message.SentAt.AsTime())
	if err != nil {
		return err
	}

	return nil
}

func ReadAllMessages(db *sqlx.DB, userId *proto.UserId) error {
	query := `UPDATE messages SET read_at = $1 WHERE recipient_id = $2`
	_, err := db.Exec(query, timestamppb.Now().AsTime(), userId.Id)
	if err != nil {
		return err
	}

	return nil
}
