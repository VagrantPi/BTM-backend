package config

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

func GetConfig(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetUserConfig")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	log.Info("GetConfig")

	data, err := repo.GetLatestConfData(repo.GetDb(c))
	if err != nil {
		log.Error("GetLatestConfData", zap.Any("err", err))
		api.ErrResponse(c, "GetLatestConfData", errors.InternalServer(error_code.ErrDiError, "GetLatestConfData").WithCause(err))
		return
	}

	api.OKResponse(c, data)
}
