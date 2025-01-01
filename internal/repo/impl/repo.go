package impl

import (
	"BTM-backend/configs"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/error_code"
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

type repository struct {
	db      *gorm.DB
	configs configs.Config
}

func NewRepository(db *gorm.DB, configs configs.Config) domain.Repository {
	return &repository{
		db:      db,
		configs: configs,
	}
}

func (repo *repository) NewTransactionBegin(ctx context.Context) (tx *gorm.DB, err error) {
	tx = repo.db.Begin().WithContext(ctx)

	if err = tx.Error; err != nil {
		err = errors.InternalServer(error_code.ErrDBError, "model.Db.Begin()").WithCause(err)
		return nil, err
	}

	return tx, nil
}

func (repo *repository) NewTxWithContext(ctx context.Context) (tx *gorm.DB, err error) {
	tx = repo.db.Begin().WithContext(ctx)

	if err = tx.Error; err != nil {
		err = errors.InternalServer(error_code.ErrDBError, "model.Db.Begin()").WithCause(err)
		return nil, err
	}

	return tx, nil
}

func (repo *repository) Close(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
	}
}

func (repo *repository) TransactionCommit(tx *gorm.DB) (err error) {
	err = tx.Commit().Error
	if err != nil {
		err = errors.InternalServer(error_code.ErrDBError, "tx.Commit()").WithCause(err)
		return
	}

	return
}

func (repo *repository) GetDb(ctx context.Context) *gorm.DB {
	return repo.db.WithContext(ctx)
}

func (repo *repository) Rollback(tx *gorm.DB) {
	tx.Rollback()
}
