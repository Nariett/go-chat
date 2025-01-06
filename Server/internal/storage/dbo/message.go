package dbo

import "time"

type Message struct {
	Id         int        `db:"id"`
	SenderId   int        `db:"sender_id"`
	ReceiverId int        `db:"receiver_id"`
	Content    string     `db:"content"`
	SentAt     time.Time  `db:"sent_at"`
	ReadAt     *time.Time `db:"read_at"`
}
