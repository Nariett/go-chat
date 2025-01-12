package activity

import (
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func UpdateLastActivity(db *sqlx.DB, id int32) error {
	query := `UPDATE activity SET date = $1 WHERE idUser = $2`
	_, err := db.Exec(query, timestamppb.Now().AsTime(), id)
	if err != nil {
		log.Fatalf("Ошибка обновления данных %v\n", err)
	}

	return nil
}
