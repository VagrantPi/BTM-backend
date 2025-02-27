package debug

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

func GetBTMChangeLogs(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetBTMChangeLogs")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	log.Info("GetBTMChangeLogs")

	data, _, err := repo.GetBTMChangeLogs(repo.GetDb(c), 100, 1)
	if err != nil {
		log.Error("GetBTMChangeLogs", zap.Any("err", err))
		api.ErrResponse(c, "GetBTMChangeLogs", errors.InternalServer(error_code.ErrDiError, "GetBTMChangeLogs").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: data,
	})
}
