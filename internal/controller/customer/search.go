package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type SearchCustomersReq struct {
	Phone string `form:"phone" binding:"required"`
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
}

type SearchCustomersRep struct {
	Code int                   `json:"code"`
	Data SearchCustomersRepItem `json:"data"`
}

type SearchCustomersRepItem struct {
	Total int               `json:"total"`
	Items []domain.Customer `json:"items"`
}

func SearchCustomers(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "SearchCustomers")
	defer func() {
		_ = log.Sync()
	}()

	req := SearchCustomersReq{}
	err := c.BindQuery(&req)
	if err != nil {
		log.Error("c.BindQuery(&req)", zap.Any("err", err))
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

	customers, total, err := repo.SearchCustomersByPhone(req.Phone, req.Limit, req.Page)
	if err != nil {
		log.Error("repo.SearchCustomersByPhone()", zap.Any("err", err))
		api.ErrResponse(c, "repo.SearchCustomersByPhone()", errors.NotFound(error_code.ErrDBError, "repo.SearchCustomersByPhone()").WithCause(err))
		return
	}

	c.JSON(
		200, SearchCustomersRep{
			Code: 20000,
			Data: SearchCustomersRepItem{
				Total: total,
				Items: customers,
			},
		},
	)
	c.Done()
}