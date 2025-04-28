package tx

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetTxsListReq struct {
	CustomerId string    `form:"customer_id"`
	Phone      string    `form:"phone"`
	DateStart  time.Time `form:"date_start"`
	DateEnd    time.Time `form:"date_end"`
	Limit      int       `form:"limit"`
	Page       int       `form:"page"`
}

type GetTxsListRep struct {
	Total int               `json:"total"`
	Items []domain.CashInTx `json:"items"`
}

func GetTxsList(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetTxsList")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := GetTxsListReq{}
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

	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	txs, total, err := repo.GetCashIns(repo.GetDb(c), req.CustomerId, req.Phone, req.DateStart, req.DateEnd, req.Limit, req.Page)
	if err != nil {
		log.Error("repo.GetCustomers()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomers()", errors.NotFound(error_code.ErrDBError, "repo.GetCustomers()").WithCause(err))
		return
	}

	api.OKResponse(c, GetTxsListRep{
		Total: total,
		Items: txs,
	})
}
