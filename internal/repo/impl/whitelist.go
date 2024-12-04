package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"

	"github.com/google/uuid"
)

func (repo *repository) CreateWhitelist(whitelist *domain.BTMWhitelist) error {
	modelWhitelist := WhitelistDomainToModel(*whitelist)
	return repo.db.Create(&modelWhitelist).Error
}

func (repo *repository) GetWhiteListByCustomerId(customerID uuid.UUID, limit int, page int) ([]domain.BTMWhitelist, int64, error) {
	var modelWhitelists []model.BTMWhitelist
	var total int64

	query := repo.db.Model(&model.BTMWhitelist{}).Where("customer_id = ?", customerID)

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

func (repo *repository) UpdateWhitelist(whitelist *domain.BTMWhitelist) error {
	modelWhitelist := WhitelistDomainToModel(*whitelist)
	return repo.db.Save(&modelWhitelist).Error
}

func (repo *repository) DeleteWhitelist(id int64) error {
	return repo.db.Delete(&model.BTMWhitelist{}, "id = ?", id).Error
}

func (repo *repository) SearchWhitelistByAddress(customerID uuid.UUID, address string, limit int, page int) ([]domain.BTMWhitelist, int64, error) {
	var modelWhitelists []model.BTMWhitelist
	var total int64

	query := repo.db.Model(&model.BTMWhitelist{}).
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
