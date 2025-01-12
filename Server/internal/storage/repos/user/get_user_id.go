package user

import "github.com/jmoiron/sqlx"

func GetUserId(db *sqlx.DB, name string) (int32, error) {
	var userId int32
	err := db.Get(&userId, "SELECT id FROM users WHERE name = $1", name)
	if err != nil {
		return -1, err
	}

	return userId, nil
}
