package riskControl

import (
	"BTM-backend/configs"
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
	RoleId        uint8           `json:"role_id"`
	RoleName      string          `json:"role_name"`
	DailyLimit    decimal.Decimal `json:"daily_limit"`
	MonthlyLimit  decimal.Decimal `json:"monthly_limit"`
	Level1        decimal.Decimal `json:"level1"`
	Level2        decimal.Decimal `json:"level2"`
	Level1Days    uint32          `json:"level1_days"`
	Level2Days    uint32          `json:"level2_days"`
	VelocityDays  uint32          `json:"velocity_days"`
	VelocityTimes uint32          `json:"velocity_times"`
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

	repo, err := di.NewRepo(configs.C.Mock)
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

	api.OKResponse(c, GetCustomerRiskControlRoleResp{
		RoleId:        userInfo.Role.Uint8(),
		RoleName:      userInfo.Role.String(),
		DailyLimit:    userInfo.DailyLimit,
		MonthlyLimit:  userInfo.MonthlyLimit,
		Level1:        userInfo.Level1,
		Level2:        userInfo.Level2,
		Level1Days:    userInfo.Level1Days,
		Level2Days:    userInfo.Level2Days,
		VelocityDays:  userInfo.VelocityDays,
		VelocityTimes: userInfo.VelocityTimes,
	})
}

type UpdateCustomerRiskControlRoleReq struct {
	CustomerId string `uri:"customer_id"`
	RoleId     uint8  `json:"role_id"`
	Reason     string `json:"reason"`
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

	if req.Reason == "" {
		log.Error("invalid reason", zap.Any("reason", req.Reason))
		api.ErrResponse(c, "invalid reason", errors.BadRequest(error_code.ErrInvalidRequest, "invalid reason").WithCause(err))
		return
	}

	repo, err := di.NewRepo(configs.C.Mock)
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

	err = repo.ChangeCustomerRole(repo.GetDb(c), operationUserInfo.Id, customerID, domain.RiskControlRole(req.RoleId), req.Reason)
	if err != nil {
		log.Error("repo.ChangeCustomerRole()", zap.Any("err", err))
		api.ErrResponse(c, "repo.ChangeCustomerRole()", errors.InternalServer(error_code.ErrDBError, "repo.ChangeCustomerRole()").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}

type ResetCustomerRiskControlRoleReq struct {
	CustomerId string `uri:"customer_id"`
}

func ResetCustomerRiskControlRole(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "ResetCustomerRiskControlRole")
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

	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	// 取得後台更改人
	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	err = repo.ResetCustomerRole(repo.GetDb(c), operationUserInfo.Id, customerID)
	if err != nil {
		log.Error("repo.ResetCustomerRole()", zap.Any("err", err))
		api.ErrResponse(c, "repo.ResetCustomerRole()", errors.InternalServer(error_code.ErrDBError, "repo.ResetCustomerRole()").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}

type ResetCustomerRiskControlRoleBatchReq struct {
	CustomerIds []string `json:"customer_ids"`
}

func ResetCustomerRiskControlRoleBatch(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "ResetCustomerRiskControlRoleBatch")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := ResetCustomerRiskControlRoleBatchReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Error("c.ShouldBindJSON(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindJSON(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindJSON(&req)").WithCause(err))
		return
	}

	if len(req.CustomerIds) == 0 {
		log.Error("empty customer_ids")
		api.ErrResponse(c, "empty customer_ids", errors.BadRequest(error_code.ErrInvalidRequest, "empty customer_ids"))
		return
	}

	customerIds := []uuid.UUID{}
	for _, item := range req.CustomerIds {
		customerID, err := uuid.Parse(item)
		if err != nil {
			log.Error("uuid.Parse(req.CustomerID)", zap.Any("err", err))
			api.ErrResponse(c, "uuid.Parse(req.CustomerID)", errors.BadRequest(error_code.ErrInvalidRequest, "uuid.Parse(req.CustomerID)").WithCause(err))
			return
		}
		customerIds = append(customerIds, customerID)
	}

	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	// 取得後台更改人
	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	for _, customerId := range customerIds {
		err = repo.ResetCustomerRole(repo.GetDb(c), operationUserInfo.Id, customerId)
		if err != nil {
			log.Error("repo.ResetCustomerRole()", zap.Any("err", err))
			api.ErrResponse(c, "repo.ResetCustomerRole()", errors.InternalServer(error_code.ErrDBError, "repo.ResetCustomerRole()").WithCause(err))
			return
		}
	}

	api.OKResponse(c, nil)
}
