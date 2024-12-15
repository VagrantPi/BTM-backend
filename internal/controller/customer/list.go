package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetCustomersListReq struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

type GetCustomersListRepItem struct {
	Total int               `json:"total"`
	Items []domain.Customer `json:"items"`
}

func GetCustomersList(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetCustomersList")
	defer func() {
		_ = log.Sync()
	}()

	req := GetCustomersListReq{}
	err := c.BindQuery(&req)
	if err != nil {
		log.Error("c.BindQuery(&req)", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery(&req)").WithCause(err))
		return
	}

	fmt.Printf("req: %+v\n", req)

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	customers, total, err := repo.GetCustomers(req.Limit, req.Page)
	if err != nil {
		log.Error("repo.GetCustomers()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomers()", errors.NotFound(error_code.ErrDBError, "repo.GetCustomers()").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: GetCustomersListRepItem{
			Total: total,
			Items: customers,
		},
	})
	c.Done()
}
