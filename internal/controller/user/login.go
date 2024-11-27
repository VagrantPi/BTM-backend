package user

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginBTMAdminReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginBTMAdminResp struct {
	Code int64                 `json:"code"`
	Data LoginBTMAdminRespItem `json:"data"`
}

type LoginBTMAdminRespItem struct {
	Token string `json:"token"`
}

func LoginBTMAdmin(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "LoginBTMAdmin")
	defer func() {
		_ = log.Sync()
	}()

	var req LoginBTMAdminReq
	err := c.BindJSON(&req)
	if err != nil {
		log.Error("c.BindJSON(&req)", zap.Any("err", err))
		_ = c.Error(err)
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		_ = c.Error(err)
		return
	}
	user, err := repo.GetBTMUserByAccount(req.Username)
	if err != nil {
		log.Error("GetBTMUserByAccount", zap.Any("err", err))
		_ = c.Error(err)
		return
	}

	if user == nil {
		log.Error("user not found")
		_ = c.Error(fmt.Errorf("login failed"))
		return
	}

	if tools.CheckPassword(user.Password, req.Password) {
		log.Error("password not match")
		_ = c.Error(fmt.Errorf("login failed"))
		return
	}

	token, err := tools.GenerateJWT(domain.UserJwt{
		Account: req.Username,
		Role:    user.Roles,
	}, configs.C.JWT.Secret)
	if err != nil {
		log.Error("GenerateJWT", zap.Any("err", err))
		_ = c.Error(err)
		return
	}

	c.JSON(
		200, LoginBTMAdminResp{
			Code: 20000,
			Data: LoginBTMAdminRespItem{
				Token: token,
			},
		},
	)
}
