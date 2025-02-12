package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SearchCustomersReq struct {
	// TODO: CICD 還沒上，這邊先魔改一版 phone 也可以查 customer_id
	Phone string `form:"phone" binding:"required"`
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
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
	c.Set("log", log)

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

	// TODO: 等 CICD 處理完，這邊需要重構或修改變數命名
	cid, _ := uuid.Parse(req.Phone)
	phone, _ := strconv.Atoi(req.Phone)
	_uuid := ""
	if cid != uuid.Nil || phone == 0 {
		// 現在 req.Phone 代表 customer_id
		_uuid = req.Phone
	}
	var customers []domain.Customer
	var total int
	if _uuid != "" {
		customers, total, err = repo.SearchCustomersByCustomerId(repo.GetDb(c), req.Phone, req.Limit, req.Page)
		if err != nil {
			log.Error("repo.SearchCustomersByCustomerId", zap.Any("err", err))
			api.ErrResponse(c, "repo.SearchCustomersByCustomerId", errors.NotFound(error_code.ErrDBError, "repo.SearchCustomersByCustomerId()").WithCause(err))
			return
		}
	} else {
		customers, total, err = repo.SearchCustomersByPhone(repo.GetDb(c), req.Phone, req.Limit, req.Page)
		if err != nil {
			log.Error("repo.SearchCustomersByPhone()", zap.Any("err", err))
			api.ErrResponse(c, "repo.SearchCustomersByPhone()", errors.NotFound(error_code.ErrDBError, "repo.SearchCustomersByPhone()").WithCause(err))
			return
		}
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: SearchCustomersRepItem{
			Total: total,
			Items: customers,
		},
	})
	c.Done()
}
