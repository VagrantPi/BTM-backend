package user

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

func LogoutBTMAdmin(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "LogoutBTMAdmin")
	defer func() {
		_ = log.Sync()
	}()

	log.Info("LogoutBTMAdmin")
	userInfo, _ := c.Get("userInfo")

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}
	err = repo.DeleteLastLoginToken(userInfo.(domain.UserJwt).Id)
	if err != nil {
		log.Error("DeleteLastLoginToken", zap.Any("err", err))
		api.ErrResponse(c, "DeleteLastLoginToken", errors.InternalServer(error_code.ErrDBError, "DeleteLastLoginToken").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
	})
}
