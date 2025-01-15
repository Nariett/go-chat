package message

import (
	proto "github.com/Nariett/go-chat/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *store) UpdateAllMessageReadStatus(userId *proto.UserId) error {
	query := `UPDATE messages SET read_at = $1 WHERE recipient_id = $2`
	_, err := s.db.Exec(query, timestamppb.Now().AsTime(), userId.Id)
	if err != nil {
		return err
	}

	return nil
}
