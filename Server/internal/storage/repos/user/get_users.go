package user

import "github.com/jmoiron/sqlx"

func GetUsers(db *sqlx.DB) ([]string, error) {
	var usernames []string
	err := db.Select(&usernames, "SELECT name FROM users")
	if err != nil {
		return nil, err
	}

	return usernames, nil
}
