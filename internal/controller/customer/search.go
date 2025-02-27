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
	Query string `form:"query" binding:"required"`
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
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
	_uuid := ""
	if cid != uuid.Nil || phone == 0 {
		// 現在搜尋條件是customer_id
		_uuid = req.Query
	}
	var customers []domain.CustomerWithWhiteListCreated
	var total int
	if _uuid != "" {
		customers, total, err = repo.SearchCustomersByCustomerId(repo.GetDb(c), _uuid, req.Limit, req.Page)
		if err != nil {
			log.Error("repo.SearchCustomersByCustomerId", zap.Any("err", err))
			api.ErrResponse(c, "repo.SearchCustomersByCustomerId", errors.NotFound(error_code.ErrDBError, "repo.SearchCustomersByCustomerId()").WithCause(err))
			return
		}
	} else {
		customers, total, err = repo.SearchCustomersByPhone(repo.GetDb(c), req.Query, req.Limit, req.Page)
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

type SearchCustomersByAddressReq struct {
	Address string `uri:"address" binding:"required"`
	Limit   int    `form:"limit"`
	Page    int    `form:"page"`
}

type SearchCustomersByAddressRepItem struct {
	Total int                                   `json:"total"`
	Items []domain.CustomerWithWhiteListCreated `json:"items"`
}

func SearchCustomersByAddress(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "SearchCustomersByAddress")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := SearchCustomersByAddressReq{}
	err := c.ShouldBindUri(&req)
	if err != nil {
		log.Error("c.ShouldBindUri(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindUri(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindUri(&req)").WithCause(err))
		return
	}

	err = c.BindQuery(&req)
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

	customers, total, err := repo.SearchCustomersByAddress(repo.GetDb(c), req.Address, req.Limit, req.Page)
	if err != nil {
		log.Error("repo.SearchCustomersByAddress", zap.Any("err", err))
		api.ErrResponse(c, "repo.SearchCustomersByAddress", errors.NotFound(error_code.ErrDBError, "repo.SearchCustomersByAddress()").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: SearchCustomersByAddressRepItem{
			Total: total,
			Items: customers,
		},
	})
	c.Done()
}

type SearchCustomersByWhitelistCreatedAtReq struct {
	DateStart time.Time `form:"date_start" binding:"required"`
	DateEnd   time.Time `form:"date_end" binding:"required"`
	Limit     int       `form:"limit"`
	Page      int       `form:"page"`
}

type SearchCustomersByWhitelistCreatedAtRepItem struct {
	Total int                                   `json:"total"`
	Items []domain.CustomerWithWhiteListCreated `json:"items"`
}

func SearchCustomersByWhitelistCreatedAt(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "SearchCustomersByWhitelistCreatedAt")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := SearchCustomersByWhitelistCreatedAtReq{}
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

	var customers []domain.CustomerWithWhiteListCreated
	var total int

	customers, total, err = repo.SearchCustomersByWhitelistCreatedAt(repo.GetDb(c), req.DateStart, req.DateEnd, req.Limit, req.Page)
	if err != nil {
		log.Error("repo.SearchCustomersByWhitelistCreatedAt", zap.Any("err", err))
		api.ErrResponse(c, "repo.SearchCustomersByWhitelistCreatedAt", errors.NotFound(error_code.ErrDBError, "repo.SearchCustomersByWhitelistCreatedAtByCustomerId()").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: SearchCustomersByWhitelistCreatedAtRepItem{
			Total: total,
			Items: customers,
		},
	})
	c.Done()
}
