package impl

import "BTM-backend/internal/repo/model"

func (repo *repository) GetBTMUserByAccount(account string) (*model.BTMUser, error) {
	info := model.BTMUser{}
	if err := repo.db.Where("account = ?", account).Find(&info).Error; err != nil {
		return nil, err
	}
	return &info, nil
}
