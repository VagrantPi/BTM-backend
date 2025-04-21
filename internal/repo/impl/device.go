package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) GetDeviceAll(db *gorm.DB) ([]domain.Device, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var devices []model.Device
	if err := db.Find(&devices).Error; err != nil {
		return nil, err
	}

	var domainDevices []domain.Device
	for _, device := range devices {
		domainDevices = append(domainDevices, DeviceModelToDomain(device))
	}
	return domainDevices, nil
}

func DeviceModelToDomain(device model.Device) domain.Device {
	return domain.Device{
		DeviceID: device.DeviceID,
		Paired:   device.Paired,
		Display:  device.Display,
		Created:  device.Created,
		Name:     device.Name,
	}
}
