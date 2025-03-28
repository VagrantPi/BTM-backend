package user

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetBTMUserInfoRespItem struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

func GetBTMUserInfo(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetBTMUserInfo")
	defer func() {
		_ = log.Sync()
	}()

	userInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("c.Get(userInfo) parse error")
		api.ErrResponse(c, "c.Get(userInfo) parse error", errors.BadRequest(error_code.ErrForbidden, "c.Get(userInfo) parse error"))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	roleAny, ok := c.Get("role")
	if !ok {
		roleAny, err = repo.GetRawRoleById(repo.GetDb(c), userInfo.Role)
		if err != nil {
			log.Error("GetRawRoleById error", zap.Any("err", err))
			api.ErrResponse(c, "GetRawRoleById error", errors.InternalServer(error_code.ErrInternalError, "GetRawRoleById error").WithCause(err))
			return
		}
	}
	role, ok := roleAny.(domain.BTMRole)
	if !ok || (ok && role.ID == 0) {
		log.Error("role not found")
		api.ErrResponse(c, "role not found", errors.NotFound(error_code.ErrDBError, "role not found"))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: GetBTMUserInfoRespItem{
			Name:  userInfo.Account,
			Roles: []string{role.RoleName},
		},
	})
	c.Done()

}
