package tests

import (
	"Server/config"
	conn "Server/internal/storage"
	"Server/internal/storage/repos/user"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserId(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		testUserName := "Иван"
		expectedUserId := int32(1)

		cnf := config.NewTestConfig()
		connection, err := conn.CreatePostgresConnection(cnf)
		db := user.NewStore(connection)

		resultUserId, err := db.GetUserId(testUserName)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, expectedUserId, resultUserId, "userId should be equal")
	})

	t.Run("failure", func(t *testing.T) {

		testUserName := "Коля"
		expectedError := sql.ErrNoRows
		expectedUserId := int32(-1)

		cnf := config.NewTestConfig()
		connection, err := conn.CreatePostgresConnection(cnf)
		db := user.NewStore(connection)

		resultUserId, err := db.GetUserId(testUserName)
		assert.ErrorIs(t, err, expectedError, "expected error for unknown user")
		assert.Equal(t, expectedUserId, resultUserId, "userId should be -1 for unknown user")
	})
}
