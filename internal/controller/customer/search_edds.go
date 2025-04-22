package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type SearchEddUsersReq struct {
	CustomerId string    `form:"customer_id"`
	Phone      string    `form:"phone"`
	Name       string    `form:"name"`
	EddStartAt time.Time `form:"edd_start_at"`
	EddEndAt   time.Time `form:"edd_end_at"`
	Limit      int       `form:"limit" binding:"required"`
	Page       int       `form:"page" binding:"required"`
}

type SearchEddUsersRep struct {
	Total int64                        `json:"total"`
	Items []domain.CustomerWithEddInfo `json:"items"`
}

func SearchEddUsers(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "SearchEddUsers")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := SearchEddUsersReq{}
	err := c.BindQuery(&req)
	if err != nil {
		log.Error("c.BindQuery(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery(&req)").WithCause(err))
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

	customers, total, err := repo.SearchEddUsers(repo.GetDb(c), req.CustomerId, req.Phone, req.Name, req.EddStartAt, req.EddEndAt, req.Limit, req.Page)
	if err != nil {
		log.Error("repo.SearchEddUsersByCustomerId", zap.Any("err", err))
		api.ErrResponse(c, "repo.SearchEddUsersByCustomerId", errors.NotFound(error_code.ErrDBError, "repo.SearchEddUsersByCustomerId()").WithCause(err))
		return
	}

	resp := SearchEddUsersRep{
		Total: total,
		Items: customers,
	}

	for i, v := range resp.Items {
		resp.Items[i].Phone = tools.MaskPhone(v.Phone)
		resp.Items[i].Name = tools.MaskName(v.Name)
	}

	api.OKResponse(c, resp)
}
