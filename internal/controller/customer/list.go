package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetCustomersListReq struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

type GetCustomersListRep struct {
	Code int                     `json:"code"`
	Data GetCustomersListRepItem `json:"data"`
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
		_ = c.Error(err)
		return
	}

	fmt.Printf("req: %+v\n", req)

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		_ = c.Error(err)
		return
	}

	customers, total, err := repo.GetCustomers(req.Limit, req.Page)
	if err != nil {
		log.Error("repo.GetCustomers()", zap.Any("err", err))
		_ = c.Error(err)
		return
	}

	c.JSON(
		200, GetCustomersListRep{
			Code: 20000,
			Data: GetCustomersListRepItem{
				Total: total,
				Items: customers,
			},
		},
	)
	c.Done()
}
