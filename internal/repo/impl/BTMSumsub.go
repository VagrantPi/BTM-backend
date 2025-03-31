package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *repository) UpsertBTMSumsub(db *gorm.DB, btmsumsub domain.BTMSumsub) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	item := BTMSumsubDomainToModel(btmsumsub)
	doUpdateItem := map[string]interface{}{
		"updated_at":           time.Now(),
		"info_hash":            item.InfoHash,
		"email_hash":           item.EmailHash,
		"phone":                item.Phone,
		"inspection_id":        item.InspectionId,
		"id_card_front_img_id": item.IdCardFrontImgId,
		"id_card_back_img_id":  item.IdCardBackImgId,
		"selfie_img_id":        item.SelfieImgId,
		"name":                 item.Name,
		"status":               item.Status,
	}

	if item.BanExpireDate.Valid {
		doUpdateItem["ban_expire_date"] = item.BanExpireDate
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id_number"}},
		DoUpdates: clause.Assignments(doUpdateItem),
	}).Create(&item).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) GetBTMSumsub(db *gorm.DB, customerId string) (*domain.BTMSumsub, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var btmsumsub model.BTMSumsub
	err := db.Where("customer_id = ?", customerId).First(&btmsumsub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		err = errors.InternalServer(error_code.ErrBTMSumsubGetItem, "GetBTMSumsub err").WithCause(err).
			WithMetadata(map[string]string{
				"customerId": customerId,
			})
		return nil, err
	}
	res := BTMSumsubModelToDomain(btmsumsub)
	return &res, nil
}

func (repo *repository) UpdateBTMSumsubBanExpireDate(db *gorm.DB, customerId string, banExpireDate int64) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	return db.Model(&model.BTMSumsub{}).Where("customer_id = ?", customerId).Update("ban_expire_date", banExpireDate).Error
}

func (repo *repository) DeleteBTMSumsub(db *gorm.DB, customerId string) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	return db.Unscoped().Delete(&model.BTMSumsub{}, "customer_id = ?", customerId).Error
}

func (repo *repository) GetUnCompletedSumsubCustomerIds(db *gorm.DB) ([]string, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var ids []string
	err := db.Model(&model.Customer{}).
		Select("customers.id").
		Joins("LEFT JOIN btm_sumsubs ON btm_sumsubs.customer_id = customers.id::text").
		Where("customers.phone != '' AND btm_sumsubs.status IS DISTINCT FROM 'GREEN'").
		Find(&ids).Error
	if err != nil {
		err = errors.InternalServer(error_code.ErrBTMSumsubGetItem, "GetBTMSumsub err").WithCause(err).
			WithMetadata(map[string]string{
				"status": "uncompleted",
			})
		return nil, err
	}
	return ids, nil
}

func BTMSumsubDomainToModel(itme domain.BTMSumsub) model.BTMSumsub {
	return model.BTMSumsub{
		CustomerId:       itme.CustomerId,
		ApplicantId:      itme.ApplicantId,
		InfoHash:         itme.InfoHash,
		IdNumber:         itme.IdNumber,
		BanExpireDate:    itme.BanExpireDate,
		EmailHash:        itme.EmailHash,
		Phone:            itme.Phone,
		InspectionId:     itme.InspectionId,
		IdCardFrontImgId: itme.IdCardFrontImgId,
		IdCardBackImgId:  itme.IdCardBackImgId,
		SelfieImgId:      itme.SelfieImgId,
		Name:             itme.Name,
		Status:           itme.Status,
	}
}

func BTMSumsubModelToDomain(itme model.BTMSumsub) domain.BTMSumsub {
	return domain.BTMSumsub{
		CustomerId:       itme.CustomerId,
		ApplicantId:      itme.ApplicantId,
		InfoHash:         itme.InfoHash,
		IdNumber:         itme.IdNumber,
		BanExpireDate:    itme.BanExpireDate,
		EmailHash:        itme.EmailHash,
		Phone:            itme.Phone,
		InspectionId:     itme.InspectionId,
		IdCardFrontImgId: itme.IdCardFrontImgId,
		IdCardBackImgId:  itme.IdCardBackImgId,
		SelfieImgId:      itme.SelfieImgId,
		Name:             itme.Name,
		Status:           itme.Status,
	}
}
