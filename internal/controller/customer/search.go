package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SearchCustomersReq struct {
	Query              string              `form:"query"`
	Address            string              `form:"address"`
	WhiteListDateStart time.Time           `form:"white_list_date_start"`
	WhiteListDateEnd   time.Time           `form:"white_list_date_end"`
	CustomerDateStart  time.Time           `form:"customer_date_start"`
	CustomerDateEnd    time.Time           `form:"customer_date_end"`
	CustomerType       domain.CustomerType `form:"customer_type"`
	Limit              int                 `form:"limit"`
	Page               int                 `form:"page"`
}

type SearchCustomersRepItem struct {
	Total int                                   `json:"total"`
	Items []domain.CustomerWithWhiteListCreated `json:"items"`
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

	cid, _ := uuid.Parse(req.Query)
	phone, _ := strconv.Atoi(req.Query)
	_customerId := ""
	if cid != uuid.Nil || phone == 0 {
		// 現在搜尋條件是customer_id
		_customerId = req.Query
		req.Query = ""
	}
	var customers []domain.CustomerWithWhiteListCreated
	var total int
	customers, total, err = repo.SearchCustomers(repo.GetDb(c), req.Query, _customerId, req.Address,
		req.WhiteListDateStart, req.WhiteListDateEnd, req.CustomerDateStart, req.CustomerDateEnd,
		req.CustomerType,
		req.Limit, req.Page)
	if err != nil {
		log.Error("repo.SearchCustomersByCustomerId", zap.Any("err", err))
		api.ErrResponse(c, "repo.SearchCustomersByCustomerId", errors.NotFound(error_code.ErrDBError, "repo.SearchCustomersByCustomerId()").WithCause(err))
		return
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
