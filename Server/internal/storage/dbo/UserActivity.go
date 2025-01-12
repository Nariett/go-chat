package dbo

import (
	"time"
)

type UserActivity struct {
	Name string    `db:"name"`
	Date time.Time `db:"date"`
}
