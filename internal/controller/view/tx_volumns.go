package view

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

type GetTxVolumnsReq struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type GetTxVolumnsRep struct {
	Items []domain.BTMDailyDeviceIncome `json:"items"`
}

func GetTxVolumns(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetTxVolumns")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := GetTxVolumnsReq{}
	err := c.BindQuery(&req)
	if err != nil {
		log.Error("c.BindQuery(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery(&req)").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	txs, err := repo.FetchByStatDate(repo.GetDb(c), req.StartDate, req.EndDate)
	if err != nil {
		log.Error("repo.GetCustomers()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomers()", errors.NotFound(error_code.ErrDBError, "repo.GetCustomers()").WithCause(err))
		return
	}

	api.OKResponse(c, GetTxVolumnsRep{
		Items: txs,
	})
}
