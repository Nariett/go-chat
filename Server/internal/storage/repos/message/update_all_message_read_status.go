package message

import (
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UpdateAllMessageReadStatus(db *sqlx.DB, userId *proto.UserId) error {
	query := `UPDATE messages SET read_at = $1 WHERE recipient_id = $2`
	_, err := db.Exec(query, timestamppb.Now().AsTime(), userId.Id)
	if err != nil {
		return err
	}

	return nil
}
