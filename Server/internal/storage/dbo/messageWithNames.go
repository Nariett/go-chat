package dbo

import "time"

type MessageWithNames struct {
	ID            int        `db:"id"`
	SenderName    string     `db:"sender_name"`
	RecipientName string     `db:"recipient_name"`
	Content       string     `db:"content"`
	SentAt        time.Time  `db:"sent_at"`
	ReadAt        *time.Time `db:"read_at"`
}
