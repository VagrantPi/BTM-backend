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

func GetBTMUserRoleRoutes(c *gin.Context) {
	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: domain.DefaultRoleRaw,
	})
	c.Done()
}

func GetBTMUserRoles(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetBTMUserRoleRoutes")
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

	list, err := repo.GetRawRoles(repo.GetDb(c))
	if err != nil {
		log.Error("GetRawRoles", zap.Any("err", err))
		api.ErrResponse(c, "GetRawRoles", errors.InternalServer(error_code.ErrDiError, "GetRawRoles").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: list,
	})
	c.Done()

}
