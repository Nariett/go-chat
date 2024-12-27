package database

import (
	"Server/config"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Nariett/go-chat/Proto"
	"github.com/lib/pq"
)

type User struct {
	name     string
	password string
}

func initDB() *sql.DB {
	cfg := config.LoadConfigWithPath("../../.env")
	connStr := cfg.BuildConnStr()
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки подключения: %s", err)
	}
	return db
}

func TestPostgreSqlConnection(t *testing.T) {
	var version string

	db := initDB()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		t.Fatalf("Ошибка подключения:%s", err)
	}

	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		t.Fatalf("Ошибка получения данных:%s", err)
	}
	fmt.Printf("Версия PostgreSQL: %s", version)
}

func TestShowAllUsersData(t *testing.T) {
	db := initDB()
	defer db.Close()

	rows, err := db.Query(`	SELECT users.id, users.name, activity.date 
    FROM users 
    JOIN activity ON users.id = activity.idUser`)
	if err != nil {
		log.Fatalf("Ошибка вывода данных %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int
			name string
			date time.Time
		)
		if err := rows.Scan(&id, &name, &date); err != nil {
			log.Fatalf("Ошибка при чтении строки: %v", err)
		}
		formattedDate := date.Format("15:04 02.01.2006")
		fmt.Printf("|ID: %d\t|Имя: %s\t|Дата: %s|\n", id, name, formattedDate)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Ошибка в процессе итерации: %v", err)
	}
}

func TestRegUser(t *testing.T) {
	user := Proto.UserData{Name: "testName", Password: "testPassword"}

	db := initDB()
	defer db.Close()

	result, err := RegUser(db, &user)
	if err != nil {
		t.Fatalf("Ошибка запроса:%s", err)
	}
	fmt.Printf("%s\n", result.Message)
}

func TestAuthUser(t *testing.T) {

	user := Proto.UserData{Name: "testName", Password: "testPassword"}

	db := initDB()
	defer db.Close()

	result, err := AuthUser(db, &user)
	if err != nil {
		t.Fatalf("Ошибка запроса:%s", err)
	}
	fmt.Println(result.Message)
}

func TestInsertTransactions(t *testing.T) {
	var userId int
	user := User{"testName", "testPassword"}

	db := initDB()
	defer db.Close()

	tr, err := db.Begin()
	if err != nil {
		t.Fatalf("Ошибка транзакции:%s", err)
	}

	err = tr.QueryRow("INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id", user.name, user.password).Scan(&userId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Printf("Данный ник уже занят: %v\n", err)
			_ = tr.Rollback()
			t.Fatalf("Данный пользователь уже занят:%s", err)

		}
		_ = tr.Rollback()
		t.Fatalf("Ошибка добавления пользователя :%s", err)

	}

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tr.Exec("INSERT INTO activity (idUser,date) VALUES ((SELECT id FROM users WHERE name = $1),$2)", user.name, time.Now().In(location))
	if err != nil {
		_ = tr.Rollback()
		t.Fatalf("Ошибка добавлению записи в таблицу activity:%s", err)
	}

	err = tr.Commit()
	if err != nil {
		t.Fatalf("Ошибка выполнения транзакции:%s", err)

	}

	log.Printf("Добавлен новый пользователь: id: %d, name: %s, password: %s\n", userId, user.name, user.password)
}

func TestDeleteTransactions(t *testing.T) {
	user := User{"testName", "testPassword"}

	db := initDB()
	defer db.Close()

	tr, err := db.Begin()
	if err != nil {
		t.Fatalf("Ошибка транзакции:%s", err)
	}

	_, err = tr.Exec("DELETE FROM activity WHERE idUser = (SELECT id FROM users WHERE name = $1 AND password = $2)", user.name, user.password)
	if err != nil {
		_ = tr.Rollback()
		t.Fatalf("Ошибка удаления таблицы activity:%s", err)
	}

	_, err = tr.Exec("DELETE FROM users WHERE name = $1 AND password = $2", user.name, user.password)
	if err != nil {
		_ = tr.Rollback()
		t.Fatalf("Ошибка удаления таблицы users:%s", err)
	}

	err = tr.Commit()
	if err != nil {
		t.Fatalf("Ошибка выполнения транзакции:%s", err)
	}

	log.Printf("Пользователь %s %s был удален из таблиц", user.name, user.password)
}
