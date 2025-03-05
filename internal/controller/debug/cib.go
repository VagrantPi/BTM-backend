package debug

import (
	"BTM-backend/internal/cronjob"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

func DownlaodCIB(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "DownlaodCIB")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	log.Info("DownlaodCIB")

	err := cronjob.DownlaodCIBAndUpsert()
	if err != nil {
		log.Error("DownlaodCIB", zap.Any("err", err))
		api.ErrResponse(c, "DownlaodCIB", errors.InternalServer(error_code.ErrDiError, "DownlaodCIB").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: nil,
	})
}
