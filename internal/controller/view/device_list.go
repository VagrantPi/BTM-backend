package view

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetDeviceListRep struct {
	Items []domain.Device `json:"items"`
}

func GetDeviceList(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetDeviceList")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	cacheKey := "system"
	resWithCache, ok := domain.GetTTLMap[domain.DeviceList](&domain.TTLDeviceListMap, cacheKey)
	if !ok {
		res, err := repo.GetDeviceAll(repo.GetDb(c))
		if err != nil {
			log.Error("repo.GetDeviceAll()", zap.Any("err", err))
			api.ErrResponse(c, "repo.GetDeviceAll()", errors.InternalServer(error_code.ErrDBError, "repo.GetDeviceAll()").WithCause(err))
			return
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
		log.Error("resWithCache is nil")
		api.ErrResponse(c, "resWithCache is nil", errors.InternalServer(error_code.ErrSumsubApiError, "resWithCache is nil"))
		return
	}

	api.OKResponse(c, GetDeviceListRep{
		Items: resWithCache.DeviceList,
	})
}
