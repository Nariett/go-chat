package activity

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func (s *store) UpdateLastActivity(id int32) error {
	query := `UPDATE activity SET date = $1 WHERE idUser = $2`
	_, err := s.db.Exec(query, timestamppb.Now().AsTime(), id)
	if err != nil {
		log.Fatalf("Ошибка обновления данных %v\n", err)
	}

	return nil
}
