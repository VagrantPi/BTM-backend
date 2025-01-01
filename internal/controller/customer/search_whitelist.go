package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SearchWhitelistReq struct {
	CustomerID string `form:"customer_id" binding:"required"`
	Address    string `form:"address" binding:"required"`
	Limit      int    `form:"limit"`
	Page       int    `form:"page"`
}

type SearchWhitelistData struct {
	Total int64                 `json:"total"`
	Items []domain.BTMWhitelist `json:"items"`
}

func SearchWhitelist(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "SearchWhitelist")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := SearchWhitelistReq{}
	err := c.BindQuery(&req)
	if err != nil {
		log.Error("c.BindQuery(&req)", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery(&req)").WithCause(err))
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		log.Error("uuid.Parse(req.CustomerID)", zap.Any("err", err))
		api.ErrResponse(c, "uuid.Parse(req.CustomerID)", errors.BadRequest(error_code.ErrInvalidRequest, "uuid.Parse(req.CustomerID)").WithCause(err))
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	customer, err := repo.GetCustomerById(repo.GetDb(c), customerID)
	if err != nil {
		log.Error("repo.GetCustomerById", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomerById", errors.NotFound(error_code.ErrDBError, "repo.GetCustomerById").WithCause(err))
		return
	}

	if customer == nil {
		log.Error("customer not found")
		api.ErrResponse(c, "customer not found", errors.NotFound(error_code.ErrDBError, "customer not found"))
		return
	}

	whitelists, total, err := repo.SearchWhitelistByAddress(repo.GetDb(c), customerID, req.Address, req.Limit, req.Page)
	if err != nil {
		log.Error("repo.SearchWhitelistByAddress", zap.Any("err", err))
		api.ErrResponse(c, "repo.SearchWhitelistByAddress", errors.NotFound(error_code.ErrDBError, "repo.SearchWhitelistByAddress").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: SearchWhitelistData{
			Total: total,
			Items: whitelists,
		},
	})
}
