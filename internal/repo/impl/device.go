package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) GetDeviceAll(db *gorm.DB) (map[string]domain.Device, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var devices []model.Device
	if err := db.Find(&devices).Error; err != nil {
		return nil, err
	}

	domainDevices := make(map[string]domain.Device, len(devices))
	for _, device := range devices {
		domainDevices[device.DeviceID] = DeviceModelToDomain(device)
	}
	return domainDevices, nil
}

func (repo *repository) GetDeviceAllWithCache(db *gorm.DB) (map[string]domain.Device, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	cacheKey := "system"
	resWithCache, ok := domain.GetTTLMap[domain.DeviceList](&domain.TTLDeviceListMap, cacheKey)
	if !ok {
		res, err := repo.GetDeviceAll(db)
		if err != nil {
			return nil, errors.InternalServer(error_code.ErrDBError, "repo.GetDeviceAll()").WithCause(err)
		}

		// 快取15天
		expire := time.Now().Add(15 * 24 * time.Hour).UnixNano()
		// cache
		roleWithTTL := domain.TTLMap[domain.DeviceList]{
			Cache: domain.DeviceList{
				DeviceList: res,
				Expiration: expire,
			},
			Expire: expire,
		}

		// 使用正確的結構調用 SetTTLMap
		domain.SetTTLMap[domain.DeviceList](&domain.TTLDeviceListMap, cacheKey, roleWithTTL.Cache, roleWithTTL.Expire)
		resWithCache = &roleWithTTL.Cache
	}

	if resWithCache == nil {
		return nil, errors.InternalServer(error_code.ErrSumsubApiError, "resWithCache is nil")
	}

	return resWithCache.DeviceList, nil
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
