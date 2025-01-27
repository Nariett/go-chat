package message

import (
	proto "github.com/Nariett/go-chat/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *store) UpdateAllMessagesFromUserReadStatus(unreadChat *proto.UnreadChat) error {
	query := `UPDATE messages SET read_at = $1 WHERE sender_id = $2 AND recipient_id = $3`
	_, err := s.db.Exec(query, timestamppb.Now().AsTime(), unreadChat.Recipient, unreadChat.Sender)
	if err != nil {
		return err
	}

	return nil
}
