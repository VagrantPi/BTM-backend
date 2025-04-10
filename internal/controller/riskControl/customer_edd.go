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
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type UpdateCustomerRiskControlEddUriReq struct {
	CustomerId string `uri:"customer_id"`
}

type UpdateCustomerRiskControlEddBodyReq struct {
	Level1 decimal.Decimal `json:"level1" binding:"required"`
	Level2 decimal.Decimal `json:"level2" binding:"required"`
}

func UpdateCustomerRiskControlEdd(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "UpdateCustomerRiskControlEdd")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	reqUri := UpdateCustomerRiskControlEddUriReq{}
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

	reqBody := UpdateCustomerRiskControlEddBodyReq{}
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

	err = repo.UpdateCustomerEdd(repo.GetDb(c), operationUserInfo.Id, customerID, reqBody.Level1, reqBody.Level2)
	if err != nil {
		log.Error("repo.UpdateCustomerEdd()", zap.Any("err", err))
		api.ErrResponse(c, "repo.UpdateCustomerEdd()", err)
		return
	}

	api.OKResponse(c, nil)
}
