package config

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetConfigLimitRes struct {
	ItemMaps map[string]domain.BTMRiskControlLimitSetting `json:"item_maps"`
}

func GetConfigLimit(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetConfigLimit")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	configs, err := repo.GetRiskControlLimitSetting(repo.GetDb(c))
	if err != nil {
		log.Error("repo.GetRiskControlLimitSetting", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetRiskControlLimitSetting", errors.InternalServer(error_code.ErrDiError, "repo.GetRiskControlLimitSetting").WithCause(err))
		return
	}

	resp := GetConfigLimitRes{
		ItemMaps: map[string]domain.BTMRiskControlLimitSetting{},
	}

	for _, config := range configs {
		switch config.Role {
		case domain.RiskControlRoleWhite:
			resp.ItemMaps["white"] = config
		case domain.RiskControlRoleGray:
			resp.ItemMaps["gray"] = config
		}
	}

	api.OKResponse(c, resp)
}
