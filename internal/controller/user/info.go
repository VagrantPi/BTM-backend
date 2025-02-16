package user

import (
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
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
