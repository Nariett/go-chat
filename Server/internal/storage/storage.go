package storage

import (
	"Server/internal/storage/repos/activity"
	"Server/internal/storage/repos/message"
	"Server/internal/storage/repos/user"

	"go.uber.org/fx"
)

func Construct() fx.Option {
	return fx.Provide(
		activity.NewStore,
		message.NewStore,
		user.NewStore,
	)
}
