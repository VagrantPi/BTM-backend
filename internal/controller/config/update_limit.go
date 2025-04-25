package config

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type UpdateConfigLimitReq struct {
	Role          domain.RiskControlRole `json:"role" binding:"required"`
	DailyLimit    decimal.Decimal        `json:"daily_limit" binding:"required"`
	MonthlyLimit  decimal.Decimal        `json:"monthly_limit" binding:"required"`
	Level1        decimal.Decimal        `json:"level1_volumn" binding:"required"`
	Level2        decimal.Decimal        `json:"level2_volumn" binding:"required"`
	Level1Days    uint32                 `json:"level1_days" binding:"required"`
	Level2Days    uint32                 `json:"level2_days" binding:"required"`
	VelocityDays  uint32                 `json:"velocity_days" binding:"required"`
	VelocityTimes uint32                 `json:"velocity_times" binding:"required"`
	Reason        string                 `json:"reason" binding:"required"`
}

func UpdateConfigLimit(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "UpdateConfigLimit")
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

	var req UpdateConfigLimitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("c.ShouldBindJSON", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindJSON", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindJSON").WithCause(err))
		return
	}

	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	// 撈出原始資料
	configs, err := repo.GetRiskControlLimitSetting(repo.GetDb(c))
	if err != nil {
		log.Error("repo.GetRiskControlLimitSetting", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetRiskControlLimitSetting", errors.InternalServer(error_code.ErrDiError, "repo.GetRiskControlLimitSetting").WithCause(err))
		return
	}
	checkFlag := false

	for _, config := range configs {
		// 如果沒有任何變更則返回
		if req.Role == config.Role {
			checkFlag = true
			if req.DailyLimit.String() == config.DailyLimit.String() &&
				req.MonthlyLimit.String() == config.MonthlyLimit.String() &&
				req.Level1.String() == config.Level1.String() &&
				req.Level2.String() == config.Level2.String() &&
				req.VelocityDays == config.VelocityDays &&
				req.VelocityTimes == config.VelocityTimes &&
				req.Level1Days == config.Level1Days &&
				req.Level2Days == config.Level2Days {
				log.Error("no limit change", zap.Any("req", req))
				api.ErrResponse(c, "no limit change", errors.BadRequest(error_code.ErrInvalidRequest, "no limit change"))
				return
			}

			break
		}
	}

	if !checkFlag {
		log.Error("role not valid", zap.Any("req", req))
		api.ErrResponse(c, "role not valid", errors.BadRequest(error_code.ErrInvalidRequest, "role not valid"))
		return
	}

	updateSetting := domain.BTMRiskControlLimitSetting{
		Role:          req.Role,
		DailyLimit:    req.DailyLimit,
		MonthlyLimit:  req.MonthlyLimit,
		Level1:        req.Level1,
		Level2:        req.Level2,
		Level1Days:    req.Level1Days,
		Level2Days:    req.Level2Days,
		VelocityDays:  req.VelocityDays,
		VelocityTimes: req.VelocityTimes,
	}

	// 只要有修改過限額或EDD的用戶，都不需要更新
	if err := repo.UpdateAllCustomerLimitSettingWithoutCustomized(repo.GetDb(c), operationUserInfo.Id, updateSetting, req.Reason); err != nil {
		log.Error("repo.UpdateAllCustomerLimitSettingWithoutCustomized", zap.Any("err", err))
		api.ErrResponse(c, "repo.UpdateAllCustomerLimitSettingWithoutCustomized", errors.InternalServer(error_code.ErrDBError, "repo.UpdateAllCustomerLimitSettingWithoutCustomized").WithCause(err))
		return
	}

	// 只更新沒修改過 Velocity 的用戶
	if err := repo.UpdateAllCustomerVelocitySettingWithoutCustomized(repo.GetDb(c), operationUserInfo.Id, updateSetting.VelocityDays, updateSetting.VelocityTimes, req.Reason); err != nil {
		log.Error("repo.UpdateAllCustomerVelocitySettingWithoutCustomized", zap.Any("err", err))
		api.ErrResponse(c, "repo.UpdateAllCustomerVelocitySettingWithoutCustomized", errors.InternalServer(error_code.ErrDBError, "repo.UpdateAllCustomerVelocitySettingWithoutCustomized").WithCause(err))
		return
	}

	// 更新
	if err := repo.UpdateRiskControlLimitSetting(repo.GetDb(c), operationUserInfo.Id, updateSetting, req.Reason); err != nil {
		log.Error("repo.UpdateRiskControlLimitSetting", zap.Any("err", err))
		api.ErrResponse(c, "repo.UpdateRiskControlLimitSetting", errors.InternalServer(error_code.ErrDBError, "repo.UpdateRiskControlLimitSetting").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}
