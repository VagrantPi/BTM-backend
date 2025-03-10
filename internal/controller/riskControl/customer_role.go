package riskControl

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type GetCustomerRiskControlRoleReq struct {
	CustomerId string `uri:"customer_id"`
}

type GetCustomerRiskControlRoleResp struct {
	RoleId       uint8           `json:"role_id"`
	RoleName     string          `json:"role_name"`
	DailyLimit   decimal.Decimal `json:"daily_limit"`
	MonthlyLimit decimal.Decimal `json:"monthly_limit"`
}

func GetCustomerRiskControlRole(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetCustomerRiskControlRole")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := GetCustomerRiskControlRoleReq{}
	err := c.ShouldBindUri(&req)
	if err != nil {
		log.Error("c.ShouldBindUri(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindUri(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindUri(&req)").WithCause(err))
		return
	}

	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		log.Error("uuid.Parse(req.CustomerID)", zap.Any("err", err))
		api.ErrResponse(c, "uuid.Parse(req.CustomerID)", errors.BadRequest(error_code.ErrInvalidRequest, "uuid.Parse(req.CustomerID)").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	userInfo, err := repo.GetRiskControlCustomerLimitSetting(repo.GetDb(c), customerID)
	if err != nil {
		log.Error("repo.GetCustomers()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomers()", errors.NotFound(error_code.ErrDBError, "repo.GetCustomers()").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: GetCustomerRiskControlRoleResp{
			RoleId:       userInfo.Role.Uint8(),
			RoleName:     userInfo.Role.String(),
			DailyLimit:   userInfo.DailyLimit,
			MonthlyLimit: userInfo.MonthlyLimit,
		},
	})
	c.Done()
}

type UpdateCustomerRiskControlRoleReq struct {
	CustomerId string `uri:"customer_id"`
	RoleId     uint8  `json:"role_id"`
}

func UpdateCustomerRiskControlRole(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "UpdateCustomerRiskControlRole")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := UpdateCustomerRiskControlRoleReq{}
	err := c.ShouldBindUri(&req)
	if err != nil {
		log.Error("c.ShouldBindUri(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindUri(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindUri(&req)").WithCause(err))
		return
	}

	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		log.Error("uuid.Parse(req.CustomerID)", zap.Any("err", err))
		api.ErrResponse(c, "uuid.Parse(req.CustomerID)", errors.BadRequest(error_code.ErrInvalidRequest, "uuid.Parse(req.CustomerID)").WithCause(err))
		return
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		log.Error("c.ShouldBindJSON(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindJSON(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindJSON(&req)").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	if (req.RoleId != domain.RiskControlRoleWhite.Uint8()) && (req.RoleId != domain.RiskControlRoleGray.Uint8()) && (req.RoleId != domain.RiskControlRoleBlack.Uint8()) {
		log.Error("invalid role id", zap.Any("role_id", req.RoleId))
		api.ErrResponse(c, "invalid role id", errors.BadRequest(error_code.ErrInvalidRequest, "invalid role id").WithCause(err))
		return
	}

	// 取得後台更改人
	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	err = repo.ChangeCustomerRole(repo.GetDb(c), operationUserInfo.Id, customerID, domain.RiskControlRole(req.RoleId))
	if err != nil {
		log.Error("repo.ChangeCustomerRole()", zap.Any("err", err))
		api.ErrResponse(c, "repo.ChangeCustomerRole()", errors.InternalServer(error_code.ErrDBError, "repo.ChangeCustomerRole()").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
	})
}
