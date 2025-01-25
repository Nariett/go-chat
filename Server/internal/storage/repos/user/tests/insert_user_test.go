package tests

import (
	"Server/config"
	conn "Server/internal/storage"
	"Server/internal/storage/repos/user"
	proto "github.com/Nariett/go-chat/Proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testUser := proto.UserData{Name: "testName", Password: "testPassword"}
		expectedResponse := &proto.ServerResponse{
			Success: true,
			Message: "Пользователь добавлен в базу данных",
		}

		cnf := config.NewTestConfig()
		connection, err := conn.CreatePostgresConnection(cnf)
		db := user.NewStore(connection)

		response, err := db.InsertUser(&testUser)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, expectedResponse, response, "response mismatch")
	})
	t.Run("failure", func(t *testing.T) {
		testUser := proto.UserData{Name: "Иван", Password: "123456"}
		expectedResponse := &proto.ServerResponse{
			Success: false,
			Message: "Пользователь с таким именем уже существует, повторите попытку.",
		}

		cnf := config.NewTestConfig()
		connection, err := conn.CreatePostgresConnection(cnf)
		db := user.NewStore(connection)

		response, err := db.InsertUser(&testUser)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, expectedResponse, response, "response mismatch")
	})
}
