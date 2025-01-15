package user

import (
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repo.go -destination=./mock/$GOFILE
type Store interface {
	GetUserId(name string) (int32, error)
	GetUsers() ([]string, error)
	InsertUser(user *proto.UserData) (*proto.ServerResponse, error)
	GetUserIdWithUpdateActivity(user *proto.UserData) (*proto.ServerResponse, error)
}

type store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{db: db}
}
