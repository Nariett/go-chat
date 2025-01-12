package message

import (
	"Server/internal/storage/dbo"
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetUnreadMessagesFromUser(db *sqlx.DB, userId *proto.UserId) ([]*proto.UserMessage, error) {
	var result []dbo.Message
	query := `SELECT * FROM messages WHERE sender_id = $1 AND recipient_id = $2 ORDER BY id`
	err := db.Select(&result, query, userId.Id, userId.Id)
	if err != nil {
		return nil, err
	}
	var messages []*proto.UserMessage
	for _, msg := range result {
		message := &proto.UserMessage{
			SenderId:    int32(msg.SenderId),
			RecipientId: int32(msg.ReceiverId),
			Content:     msg.Content,
			SentAt:      timestamppb.New(msg.SentAt),
		}
		messages = append(messages, message)
	}
	return messages, nil
}
