package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (repo *repository) CreateWhitelist(db *gorm.DB, whitelist *domain.BTMWhitelist) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	modelWhitelist := WhitelistDomainToModel(*whitelist)
	return db.Create(&modelWhitelist).Error
}

func (repo *repository) UpdateWhitelistSoftDelete(db *gorm.DB, whitelist *domain.BTMWhitelist) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	modelWhitelist := WhitelistDomainToModel(*whitelist)
	return db.Model(&model.BTMWhitelist{}).
		Where("customer_id = ? AND crypto_code = ? AND address =?", modelWhitelist.CustomerID, modelWhitelist.CryptoCode, modelWhitelist.Address).
		Unscoped().
		Updates(map[string]interface{}{"deleted_at": nil, "created_at": time.Now()}).Error
}

func (repo *repository) GetWhiteListByCustomerId(db *gorm.DB, customerID uuid.UUID, limit int, page int) ([]domain.BTMWhitelist, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var modelWhitelists []model.BTMWhitelist
	var total int64

	query := db.Model(&model.BTMWhitelist{}).Where("customer_id = ?", customerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&modelWhitelists).Error; err != nil {
		return nil, 0, err
	}

	whitelists := make([]domain.BTMWhitelist, len(modelWhitelists))
	for i, modelWhitelist := range modelWhitelists {
		whitelists[i] = WhitelistModelToDomain(modelWhitelist)
	}
	return whitelists, total, nil
}

func (repo *repository) GetWhiteLists(db *gorm.DB) (list []domain.BTMWhitelist, err error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var modelWhitelists []model.BTMWhitelist
	if err := db.Find(&modelWhitelists).Error; err != nil {
		return nil, err
	}

	whitelists := make([]domain.BTMWhitelist, len(modelWhitelists))
	for i, modelWhitelist := range modelWhitelists {
		whitelists[i] = WhitelistModelToDomain(modelWhitelist)
	}
	return whitelists, nil
}

func (repo *repository) CheckExistWhitelist(db *gorm.DB, customerID uuid.UUID, cryptoCode string, address string, isUnscoped bool) (bool, bool, error) {
	if db == nil {
		return false, false, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var fetched model.BTMWhitelist
	sql := db.Where("customer_id = ? AND crypto_code = ? AND address = ?", customerID, cryptoCode, address)

	var err error
	if isUnscoped {
		err = sql.Unscoped().First(&fetched).Error
	} else {
		err = sql.First(&fetched).Error
	}

	exist := !errors.Is(err, gorm.ErrRecordNotFound)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return exist, fetched.DeletedAt.Valid, err
}

func (repo *repository) UpdateWhitelist(db *gorm.DB, whitelist *domain.BTMWhitelist) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	modelWhitelist := WhitelistDomainToModel(*whitelist)
	return db.Save(&modelWhitelist).Error
}

func (repo *repository) GetWhiteListById(db *gorm.DB, id int64) (data domain.BTMWhitelist, err error) {
	if db == nil {
		return domain.BTMWhitelist{}, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var modelWhitelist model.BTMWhitelist
	if err := db.First(&modelWhitelist, "id = ?", id).Error; err != nil {
		return domain.BTMWhitelist{}, err
	}

	return WhitelistModelToDomain(modelWhitelist), nil
}

func (repo *repository) GetWhiteListByAddress(db *gorm.DB, address string) (data domain.BTMWhitelist, err error) {
	if db == nil {
		return domain.BTMWhitelist{}, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var modelWhitelist model.BTMWhitelist
	if err := db.First(&modelWhitelist, "address = ?", address).Error; err != nil {
		return domain.BTMWhitelist{}, err
	}

	return WhitelistModelToDomain(modelWhitelist), nil
}

func (repo *repository) DeleteWhitelist(db *gorm.DB, id int64) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	return db.Delete(&model.BTMWhitelist{}, "id = ?", id).Error
}

func (repo *repository) SearchWhitelistByAddress(db *gorm.DB, customerID uuid.UUID, address string, limit int, page int) ([]domain.BTMWhitelist, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var modelWhitelists []model.BTMWhitelist
	var total int64

	query := db.Model(&model.BTMWhitelist{}).
		Where("customer_id = ?", customerID).
		Where("address LIKE ?", address+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&modelWhitelists).Error; err != nil {
		return nil, 0, err
	}

	whitelists := make([]domain.BTMWhitelist, len(modelWhitelists))
	for i, modelWhitelist := range modelWhitelists {
		whitelists[i] = WhitelistModelToDomain(modelWhitelist)
	}
	return whitelists, total, nil
}

func WhitelistDomainToModel(whitelist domain.BTMWhitelist) model.BTMWhitelist {
	return model.BTMWhitelist{
		CustomerID: whitelist.CustomerID,
		CryptoCode: whitelist.CryptoCode,
		Address:    whitelist.Address,
	}
}

func WhitelistModelToDomain(whitelist model.BTMWhitelist) domain.BTMWhitelist {
	return domain.BTMWhitelist{
		ID:         uint64(whitelist.ID),
		CustomerID: whitelist.CustomerID,
		CryptoCode: whitelist.CryptoCode,
		Address:    whitelist.Address,
		CreatedAt:  whitelist.CreatedAt,
		UpdatedAt:  whitelist.UpdatedAt,
	}
}
