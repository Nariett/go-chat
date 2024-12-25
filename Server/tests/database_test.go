package tests

import (
	"Server/config"
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

// func TestRegUser(t *testing.T) {

// }

func TestAuthUser(t *testing.T) {

}

func TestPostgreSqlConnection(t *testing.T) {
	config := config.LoadConfig()
	connStr := config.BuildConnStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatalf("Ошибка подключения:%s", err)
	}

	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatalf("Ошибка вывода данных %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id   int
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("Ошибка при чтении строки: %v", err)
		}
		fmt.Printf("ID: %d, Имя: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Ошибка в процессе итерации: %v", err)
	}
}
