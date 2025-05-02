package system

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

type GetMockTxHistoryLogsReq struct {
	Limit      int       `form:"limit" binding:"required"`
	Page       int       `form:"page" binding:"required"`
	StartAt    time.Time `form:"start_at"`
	EndAt      time.Time `form:"end_at"`
	CustomerId string    `form:"customer_id"`
}

type GetMockTxHistoryLogsResp struct {
	Items []domain.BTMMockTxHistoryLog `json:"items"`
	Total int64                        `json:"total"`
}

func GetMockTxHistoryLogs(c *gin.Context) {
	req := GetMockTxHistoryLogsReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		api.ErrResponse(c, "c.ShouldBindQuery(\"req\")", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindQuery(\"req\")"))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	logs, total, err := repo.GetMockTxHistoryLogs(repo.GetDb(c), req.Limit, req.Page, req.StartAt, req.EndAt, req.CustomerId)
	if err != nil {
		api.ErrResponse(c, "repo.GetMockTxHistoryLogs()", errors.InternalServer(error_code.ErrDBError, "repo.GetMockTxHistoryLogs()").WithCause(err))
		return
	}

	api.OKResponse(c, GetMockTxHistoryLogsResp{Items: logs, Total: total})
}
