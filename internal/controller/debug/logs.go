package debug

import (
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

type GetBTMChangeLogsReq struct {
	TableName  string    `form:"table_name"`
	CustomerId string    `form:"customer_id"`
	StartAt    time.Time `form:"start_at"`
	EndAt      time.Time `form:"end_at"`
	Limit      int       `form:"limit"`
	Page       int       `form:"page"`
}

type GetBTMChangeLogsRep struct {
	Total int64                 `json:"total"`
	Items []domain.BTMChangeLog `json:"items"`
}

func GetBTMChangeLogs(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetBTMChangeLogs")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := &GetBTMChangeLogsReq{}
	if err := c.BindQuery(req); err != nil {
		log.Error("c.BindQuery()", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery()", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery()").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	log.Info("GetBTMChangeLogs")

	list, total, err := repo.GetBTMChangeLogs(repo.GetDb(c), req.TableName, req.CustomerId, req.StartAt, req.EndAt, req.Limit, req.Page)
	if err != nil {
		log.Error("GetBTMChangeLogs", zap.Any("err", err))
		api.ErrResponse(c, "GetBTMChangeLogs", errors.InternalServer(error_code.ErrDiError, "GetBTMChangeLogs").WithCause(err))
		return
	}

	api.OKResponse(c, GetBTMChangeLogsRep{
		Total: total,
		Items: list,
	})
}
