package user

import (
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"strings"

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

	token := c.GetHeader("token")
	if token == "" {
		log.Error("token is empty")
		api.ErrResponse(c, "token is empty", errors.BadRequest(error_code.ErrForbidden, "token is empty"))
		return
	}

	userInfo, err := tools.ParseToken(token)
	if err != nil {
		if strings.Contains(err.Error(), "Token is expired") {
			api.ErrResponse(c, "Token is expired", errors.BadRequest(error_code.ErrTokenExpired, "Token is expired").WithCause(err))
			return
		}
		log.Error("parseToken error", zap.Any("err", err))
		api.ErrResponse(c, "parseToken error", errors.BadRequest(error_code.ErrForbidden, "parseToken error").WithCause(err))
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
