package hotfix

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

func CustomerLimitFillDefaultConfig(c *gin.Context) {
	repo, err := di.NewRepo()
	if err != nil {
		fmt.Println("di.NewRepo()", err)
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	settings, err := repo.GetAllCustomerLimitSetting(repo.GetDb(c))
	if err != nil {
		api.ErrResponse(c, "repo.GetAllCustomerLimitSetting()", errors.InternalServer(error_code.ErrDBError, "repo.GetAllCustomerLimitSetting()").WithCause(err))
		return
	}
	defaultConfigs, err := repo.GetRiskControlLimitSetting(repo.GetDb(c))
	if err != nil {
		api.ErrResponse(c, "repo.GetRiskControlLimitSetting", errors.InternalServer(error_code.ErrDiError, "repo.GetRiskControlLimitSetting").WithCause(err))
		return
	}

	defaultConfigMap := make(map[domain.RiskControlRole]domain.BTMRiskControlLimitSetting)
	for _, defaultConfig := range defaultConfigs {
		if defaultConfig.Role == domain.RiskControlRoleWhite {
			defaultConfigMap[domain.RiskControlRoleWhite] = defaultConfig
		}
		if defaultConfig.Role == domain.RiskControlRoleGray {
			defaultConfigMap[domain.RiskControlRoleGray] = defaultConfig
		}
	}

	// TODO: hotfix 用完後應刪除，量少情況下才能呼叫
	for _, setting := range settings {
		_, err := repo.GetRiskControlCustomerLimitSetting(repo.GetDb(c), setting.CustomerId)
		if err != nil {
			api.ErrResponse(c, "repo.GetRiskControlCustomerLimitSetting", errors.InternalServer(error_code.ErrDBError, "repo.GetRiskControlCustomerLimitSetting").WithCause(err))
			return
		}
	}

	api.OKResponse(c, nil)
}
