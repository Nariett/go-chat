package message

import (
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
	"log"
)

func InsertMessage(db *sqlx.DB, message *proto.UserMessage) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Printf("Ошибка начала транзакции %v\n", err)
		return err
	}

	query := `INSERT INTO messages (sender_id, recipient_id, content, sent_at) 
			  VALUES ($1, $2, $3, $4)`
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
