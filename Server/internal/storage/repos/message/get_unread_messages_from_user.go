package message

import (
	"Server/internal/storage/dbo"
	proto "github.com/Nariett/go-chat/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *store) GetUnreadMessagesFromUser(user *proto.UnreadChat) (*proto.UserMessages, error) {
	var result []dbo.Message
	query := `SELECT * FROM messages WHERE sender_id = $1 AND recipient_id = $2 AND read_at IS NULL ORDER BY id`
	err := s.db.Select(&result, query, user.Recipient, user.Sender)
	if err != nil {
		return nil, err
	}
	var messagesArray []*proto.UserMessage
	for _, msg := range result {
		message := &proto.UserMessage{
			SenderId:    int32(msg.SenderId),
			RecipientId: int32(msg.ReceiverId),
			Content:     msg.Content,
			SentAt:      timestamppb.New(msg.SentAt),
		}
		messagesArray = append(messagesArray, message)
	}
	messages := &proto.UserMessages{
		Messages: messagesArray,
	}

	return messages, nil
}
