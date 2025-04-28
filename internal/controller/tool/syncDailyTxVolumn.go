package tool

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

type SyncDailyTxVolumnReq struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

func SyncDailyTxVolumn(c *gin.Context) {
	var req SyncDailyTxVolumnReq
	if err := c.ShouldBindQuery(&req); err != nil {
		api.ErrResponse(c, "SyncDailyTxVolumn", errors.BadRequest(error_code.ErrInvalidRequest, "SyncDailyTxVolumn").WithCause(err))
		return
	}

	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		api.ErrResponse(c, "SyncDailyTxVolumn", errors.InternalServer(error_code.ErrDiError, "SyncDailyTxVolumn").WithCause(err))
		return
	}

	err = repo.SnapshotRange(repo.GetDb(c), req.StartDate, req.EndDate)
	if err != nil {
		api.ErrResponse(c, "SyncDailyTxVolumn", errors.InternalServer(error_code.ErrDBError, "SyncDailyTxVolumn").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}
