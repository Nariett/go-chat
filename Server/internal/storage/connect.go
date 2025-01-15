package storage

import (
	"Server/config"
	"github.com/jmoiron/sqlx"
)

func CreatePostgresConnection(cfg *config.Config) (*sqlx.DB, error) {
	connStr := cfg.BuildConnStr()
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
