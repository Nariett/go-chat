package message

import (
	"Server/internal/storage/dbo"
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
)

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
