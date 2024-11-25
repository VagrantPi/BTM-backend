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

func NewRepo() (domain.Repository, error) {
	panic(wire.Build(impl.NewRepository, db.ConnectToDatabase, configs.NewConfigs))
}
