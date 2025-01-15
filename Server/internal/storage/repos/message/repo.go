package message

import (
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	GetUnreadMessagesCounter(id *proto.UserId) (*proto.UnreadMessages, error)
	GetUnreadMessagesFromUser(userId *proto.UserId) ([]*proto.UserMessage, error)
	InsertMessage(message *proto.UserMessage) error
	UpdateAllMessageReadStatus(userId *proto.UserId) error
	UpdateMessageReadStatus(message *proto.UserMessage) error
}

type store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{db: db}
}
