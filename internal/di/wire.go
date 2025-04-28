//go:build wireinject
// +build wireinject

package di

import (
	"BTM-backend/configs"
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/impl"
	"BTM-backend/third_party/db"

	"github.com/google/wire"
)

func NewRepo(isMock bool) (domain.Repository, error) {
	wire.Build(
		impl.NewRepository,
		configs.NewConfigs,
		db.ProvideDatabase,
	)
	return nil, nil // << 這是必要的，Wire 的語法規定要 return 一個空的
}
