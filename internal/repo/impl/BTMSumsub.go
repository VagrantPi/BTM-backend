package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) CreateBTMSumsub(db *gorm.DB, btmsumsub domain.BTMSumsub) error {
	item := BTMSumsubDomainToModel(btmsumsub)
	if err := db.Create(&item).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) GetBTMSumsub(db *gorm.DB, customerId string) (*domain.BTMSumsub, error) {
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

func BTMSumsubDomainToModel(itme domain.BTMSumsub) model.BTMSumsub {
	return model.BTMSumsub{
		CustomerId: itme.CustomerId,
		Info:       itme.Info,
		IdNumber:   itme.IdNumber,
	}
}

func BTMSumsubModelToDomain(itme model.BTMSumsub) domain.BTMSumsub {
	return domain.BTMSumsub{
		CustomerId: itme.CustomerId,
		Info:       itme.Info,
		IdNumber:   itme.IdNumber,
	}
}
