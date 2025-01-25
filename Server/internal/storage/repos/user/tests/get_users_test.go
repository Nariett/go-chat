package tests

import (
	mock "Server/internal/storage/repos/user/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mUser := mock.NewMockStore(ctrl)
		expectedUsers := []string{"Иван", "Мария", "Александр", "Ольга", "Дмитрий", "Анна", "Екатерина", "Михаил", "Наталья", "Сергей"}
		mUser.EXPECT().GetUsers().Return(expectedUsers, nil)

		resultUserId, err := mUser.GetUsers()
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, expectedUsers, resultUserId, "Users not retrieved correctly")
	})
}
