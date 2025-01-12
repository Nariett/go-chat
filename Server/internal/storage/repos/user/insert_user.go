package user

import (
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func InsertUser(db *sqlx.DB, user *proto.UserData) (*proto.ServerResponse, error) {
	tx, err := db.Beginx()
	if err != nil {
		return &proto.ServerResponse{
			Success: false,
			Message: "Ошибка начала транзакции",
		}, err
	}

	var existingUserId int32
	checkQuery := `SELECT id FROM users WHERE name = $1`
	err = tx.Get(&existingUserId, checkQuery, user.Name)
	if err == nil {
		log.Printf("Данный ник уже занят: %s\n", user.Name)
		_ = tx.Rollback()
		return &proto.ServerResponse{
			Success: false,
			Message: "Пользователь с таким именем уже существует, повторите попытку.",
		}, nil
	}

	var userId int32

	query := `INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id`
	err = tx.Get(&userId, query, user.Name, user.Password)
	if err != nil {
		log.Println("Ошибка транзакции")
		_ = tx.Rollback()
		return &proto.ServerResponse{
			Success: false,
			Message: "Пользователь не добавлен в базу данных, ошибка базы данных, повторите попытку",
		}, nil
	}

	query = `INSERT INTO activity (idUser,date) VALUES ($1,$2)`
	_, err = tx.Exec(query, userId, timestamppb.Now().AsTime())
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
