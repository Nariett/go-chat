package activity

import (
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	GetUsersActivityDates(_ *proto.Empty) (*proto.UserActivityDates, error)
	UpdateLastActivity(id int32) error
}

type store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{db: db}
}
