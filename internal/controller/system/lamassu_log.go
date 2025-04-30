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

type GetLamassuLogsReq struct {
	Limit   int       `form:"limit" binding:"required"`
	Page    int       `form:"page" binding:"required"`
	StartAt time.Time `form:"start_at"`
	EndAt   time.Time `form:"end_at"`
}

type GetLamassuLogsData struct {
	Total int64              `json:"total"`
	Items []domain.ServerLog `json:"items"`
}

func GetLamassuLogs(c *gin.Context) {
	req := GetLamassuLogsReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		api.ErrResponse(c, "c.ShouldBindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindQuery(&req)").WithCause(err))
		return
	}

	if req.Limit <= 0 {
		req.Limit = 100
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	repo, err := di.NewRepo()
	if err != nil {
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	logs, total, err := repo.GetServerLogs(repo.GetDb(c), req.Limit, req.Page, req.StartAt, req.EndAt)
	if err != nil {
		api.ErrResponse(c, "repo.GetServerLogs", errors.NotFound(error_code.ErrDBError, "repo.GetServerLogs").WithCause(err))
		return
	}

	api.OKResponse(c, GetLamassuLogsData{
		Total: total,
		Items: logs,
	})
}
