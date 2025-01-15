package message

import (
	proto "github.com/Nariett/go-chat/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *store) UpdateMessageReadStatus(message *proto.UserMessage) error {
	query := `UPDATE messages SET read_at = $1 WHERE content = $2 AND sent_at = $3`
	_, err := s.db.Exec(query, timestamppb.Now().AsTime(), message.Content, message.SentAt.AsTime())
	if err != nil {
		return err
	}

	return nil
}
