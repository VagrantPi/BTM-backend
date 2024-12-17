package user

import (
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

	u, exist := c.Get("userInfo")
	if !exist {
		log.Error("c.Get(userInfo) error")
		api.ErrResponse(c, "c.Get(userInfo) error", errors.BadRequest(error_code.ErrForbidden, "c.Get(userInfo) error"))
		return
	}

	userInfo, isUser := u.(domain.UserJwt)
	if !isUser {
		log.Error("c.Get(userInfo) parse error", zap.Any("u", u))
		api.ErrResponse(c, "c.Get(userInfo) parse error", errors.BadRequest(error_code.ErrForbidden, "c.Get(userInfo) parse error"))
		return
	}

	roles := tools.Role(userInfo.Role)
	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: GetBTMUserInfoRespItem{
			Name:  userInfo.Account,
			Roles: roles.ToStrings(),
		},
	})
	c.Done()

}
