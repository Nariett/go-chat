package user

import (
	"database/sql"
	"errors"
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func GetUserIdWithUpdateActivity(db *sqlx.DB, user *proto.UserData) (*proto.ServerResponse, error) {
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
