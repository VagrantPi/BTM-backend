package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"fmt"
	"time"

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

func (repo *repository) DeleteBTMCIB(db *gorm.DB, pid string) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	return db.Delete(&model.BTM_CIB{}, "pid = ?", pid).Error
}

func (repo *repository) GetBTMCIBs(db *gorm.DB, id string, limit int, page int) ([]domain.BTMCIB, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.BTM_CIB{}

	// 取得現在的中華民國年日期
	now := time.Now()
	twYear := now.Year() - 1911
	today := fmt.Sprintf("%03d%02d%02d", twYear, int(now.Month()), now.Day())

	sql := db.Model(&model.BTM_CIB{}).Where("data_type != 'D' and expire_date >= ?", today)
	if id != "" {
		sql = sql.Where("UPPER(TRIM(pid)) = UPPER(TRIM(?))", id)
	}

	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.BTMCIB, 0, len(list))
	for _, item := range list {
		resp = append(resp, BTMCIBModelToDomain(item))
	}
	return resp, int64(total), nil

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
