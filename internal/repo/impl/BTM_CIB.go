package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *repository) UpsertBTMCIB(db *gorm.DB, cib domain.BTMCIB) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	item := BTMCIBDomainToModel(cib)

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "pid"}},
		DoNothing: true,
	}).Create(&item).Error
}

func BTMCIBDomainToModel(item domain.BTMCIB) model.BTM_CIB {
	return model.BTM_CIB{
		DataType:    item.DataType,
		Pid:         item.Pid,
		WarningDate: item.WarningDate,
		ExpireDate:  item.ExpireDate,
		Issuer:      item.Issuer,
		Blank1:      item.Blank1,
		Blank2:      item.Blank2,
	}
}

func BTMCIBModelToDomain(item model.BTM_CIB) domain.BTMCIB {
	return domain.BTMCIB{
		DataType:    item.DataType,
		Pid:         item.Pid,
		WarningDate: item.WarningDate,
		ExpireDate:  item.ExpireDate,
		Issuer:      item.Issuer,
		Blank1:      item.Blank1,
		Blank2:      item.Blank2,
	}
}
