package impl

import (
	"BTM-backend/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (repo *repository) GetCustomerLimitChangeLogs(db *gorm.DB, customerID uuid.UUID, start, end time.Time, page, limit int) ([]domain.BTMRiskControlCustomerLimitSettingChange, int64, error) {
	// if db == nil {
	// 	return nil, 0, errors.BadRequest(error_code.ErrInvalidRequest, "db is nil")
	// }

	// sql := db.Model(&model.BTMChangeLog{}).Where("customer_id = ?", customerID)

	// offset := (page - 1) * limit
	// list := []model.BTMRiskControlCustomerLimitSettingChange{}

	// if !start.IsZero() && !end.IsZero() {
	// 	sql = sql.Where("created_at BETWEEN ? AND ?", start, end)
	// }

	// var total int64 = 0
	// if err := sql.Count(&total).Error; err != nil {
	// 	return nil, 0, err
	// }

	// if err := sql.
	// 	Offset(offset).
	// 	Limit(limit).
	// 	Order("created_at desc").
	// 	Find(&list).Error; err != nil {
	// 	return nil, 0, err
	// }

	// resp := make([]domain.BTMRiskControlCustomerLimitSettingChange, len(list))
	// for i, v := range list {
	// 	resp[i] = BTMRiskControlCustomerLimitSettingChangeModelToDomain(v)
	// }

	// return resp, total, nil
	return nil, 0, nil
}
