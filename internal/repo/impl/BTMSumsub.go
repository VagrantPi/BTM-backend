package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"strings"
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

func (repo *repository) GetUnCompletedSumsubCustomerIds(db *gorm.DB, force bool) ([]string, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var ids []string
	sql := db.Model(&model.Customer{}).
		Select("customers.id").
		Joins("LEFT JOIN btm_sumsubs ON btm_sumsubs.customer_id = customers.id::text")

	if force {
		sql = sql.Where("customers.phone != ''")
	} else {
		sql = sql.Where("(customers.phone != '' AND btm_sumsubs.status IS DISTINCT FROM 'GREEN') OR btm_sumsubs.updated_at < ?", time.Now().AddDate(0, 0, -10))
	}
	if err := sql.Find(&ids).Error; err != nil {
		err = errors.InternalServer(error_code.ErrBTMSumsubGetItem, "GetBTMSumsub err").WithCause(err).
			WithMetadata(map[string]string{
				"status": "uncompleted",
			})
		return nil, err
	}
	return ids, nil
}

func (repo *repository) SearchEddUsers(db *gorm.DB, customerId, phone, name string, eddStartAt, eddEndAt time.Time, limit int, page int) ([]domain.CustomerWithEddInfo, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var customers []domain.CustomerWithEddInfo
	sql := db.Model(&model.BTMSumsub{}).
		Select(
			"btm_risk_control_customer_limit_settings.edd_at",
			"btm_sumsubs.customer_id",
			"btm_sumsubs.name",
			"btm_sumsubs.phone",
			"btm_risk_control_customer_limit_settings.edd_type",
		).
		Joins("INNER JOIN btm_risk_control_customer_limit_settings ON btm_sumsubs.customer_id = btm_risk_control_customer_limit_settings.customer_id").
		Where("btm_risk_control_customer_limit_settings.is_edd = TRUE")

	switch {
	case strings.TrimSpace(customerId) != "":
		sql = sql.Where("btm_sumsubs.customer_id = ?", strings.TrimSpace(customerId))
	case strings.TrimSpace(phone) != "":
		sql = sql.Where("btm_sumsubs.phone LIKE ?", "%"+strings.TrimSpace(phone)+"%")
	case strings.TrimSpace(name) != "":
		sql = sql.Where("btm_sumsubs.name = ?", strings.TrimSpace(name))
	}

	if !eddStartAt.IsZero() && !eddEndAt.IsZero() {
		sql = sql.Where("btm_risk_control_customer_limit_settings.edd_at BETWEEN ? AND ?", eddStartAt, eddEndAt)
	}

	var total int64
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := sql.
		Limit(limit).
		Offset(offset).
		Find(&customers).Error
	if err != nil {
		err = errors.InternalServer(error_code.ErrBTMSumsubGetItem, "GetBTMSumsub err").WithCause(err).
			WithMetadata(map[string]string{
				"status": "green",
			})
		return nil, 0, err
	}

	return customers, total, nil
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
