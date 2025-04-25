package riskControl

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UpdateCustomerRiskControlVelocityUriReq struct {
	CustomerId string `uri:"customer_id"`
}

type UpdateCustomerRiskControlVelocityBodyReq struct {
	VelocityDays  uint32 `json:"velocity_days" binding:"required"`
	VelocityTimes uint32 `json:"velocity_times" binding:"required"`
	Reason        string `json:"reason" binding:"required"`
}

func UpdateCustomerRiskControlVelocity(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "UpdateCustomerRiskControlVelocity")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	reqUri := UpdateCustomerRiskControlVelocityUriReq{}
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		log.Error("c.ShouldBindUri(reqUri)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindUri(&reqUri)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindUri(&req)").WithCause(err))
		return
	}

	customerID, err := uuid.Parse(reqUri.CustomerId)
	if err != nil {
		log.Error("uuid.Parse(req.CustomerID)", zap.Any("err", err))
		api.ErrResponse(c, "uuid.Parse(req.CustomerID)", errors.BadRequest(error_code.ErrInvalidRequest, "uuid.Parse(req.CustomerID)").WithCause(err))
		return
	}

	reqBody := UpdateCustomerRiskControlVelocityBodyReq{}
	err = c.ShouldBindJSON(&reqBody)
	if err != nil {
		log.Error("c.ShouldBindJSON(reqBody)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindJSON(&reqBody)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindJSON(&reqBody)").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
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

	// 更新限額
	err = repo.UpdateCustomerVelocity(repo.GetDb(c), operationUserInfo.Id, customerID, reqBody.VelocityDays, reqBody.VelocityTimes, reqBody.Reason)
	if err != nil {
		log.Error("repo.UpdateCustomerVelocity()", zap.Any("err", err))
		api.ErrResponse(c, "repo.UpdateCustomerVelocity()", err)
		return
	}

	api.OKResponse(c, nil)
}
