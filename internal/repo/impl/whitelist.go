package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *repository) CreateWhitelist(db *gorm.DB, whitelist *domain.BTMWhitelist) error {
	modelWhitelist := WhitelistDomainToModel(*whitelist)
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "customer_id"}, {Name: "crypto_code"}, {Name: "address"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"customer_id": whitelist.CustomerID,
			"crypto_code": whitelist.CryptoCode,
			"address":     whitelist.Address,
			"deleted_at":  nil,
		}),
	}).Create(&modelWhitelist).Error
}

func (repo *repository) GetWhiteListByCustomerId(db *gorm.DB, customerID uuid.UUID, limit int, page int) ([]domain.BTMWhitelist, int64, error) {
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

func (repo *repository) UpdateWhitelist(db *gorm.DB, whitelist *domain.BTMWhitelist) error {
	modelWhitelist := WhitelistDomainToModel(*whitelist)
	return db.Save(&modelWhitelist).Error
}

func (repo *repository) DeleteWhitelist(db *gorm.DB, id int64) error {
	return db.Delete(&model.BTMWhitelist{}, "id = ?", id).Error
}

func (repo *repository) SearchWhitelistByAddress(db *gorm.DB, customerID uuid.UUID, address string, limit int, page int) ([]domain.BTMWhitelist, int64, error) {
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
