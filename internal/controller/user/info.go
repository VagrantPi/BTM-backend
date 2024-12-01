package user

import (
	"BTM-backend/configs"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetBTMUserInfoResp struct {
	Code int64                  `json:"code"`
	Data GetBTMUserInfoRespItem `json:"data"`
}

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

	userInfo, err := parseToken(token)
	if err != nil {
		log.Error("parseToken error", zap.Any("err", err))
		api.ErrResponse(c, "parseToken error", errors.BadRequest(error_code.ErrForbidden, "parseToken error").WithCause(err))
		return
	}

	roles := tools.Role(userInfo.Role)
	c.JSON(
		200, GetBTMUserInfoResp{
			Code: 20000,
			Data: GetBTMUserInfoRespItem{
				Name:  userInfo.Account,
				Roles: roles.ToStrings(),
			},
		},
	)
	c.Done()

}

func parseToken(token string) (claim domain.UserJwt, err error) {
	// var data []byte
	data, err := tools.ParseJWT(token, configs.C.JWT.Secret)
	if err != nil {
		err = errors.Unauthorized(error_code.ErrInvalidJWTParse, "ParseJWT").WithCause(err)
		return
	}

	err = json.Unmarshal(data, &claim)
	if err != nil {
		err = errors.Unauthorized(error_code.ErrInvalidJWT, "Unmarshal").WithCause(err)
		return
	}

	return
}
