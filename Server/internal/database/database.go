package database

import (
	"database/sql"
	"errors"
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

func UpdateLastActivity(db *sql.DB, user *proto.User) error {
	location, err := time.LoadLocation("Europe/Moscow")
	_, err = db.Exec("UPDATE activity SET date = $1 WHERE idUser = (SELECT id FROM users WHERE name = $2)", time.Now().In(location), user.Name)
	if err != nil {
		log.Fatalf("Обновленрия данных %v\n", err)
	}
	return nil
}
