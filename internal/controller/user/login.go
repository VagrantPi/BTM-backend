package user

import (
	"BTM-backend/configs"
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

type LoginBTMAdminReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginBTMAdminRespItem struct {
	Token string `json:"token"`
}

func LoginBTMAdmin(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "LoginBTMAdmin")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	var req LoginBTMAdminReq
	err := c.BindJSON(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}
	user, err := repo.GetBTMUserByAccount(repo.GetDb(c), req.Username)
	if err != nil {
		log.Error("GetBTMUserByAccount", zap.Any("err", err))
		api.ErrResponse(c, "GetBTMUserByAccount", errors.InternalServer(error_code.ErrDiError, "GetBTMUserByAccount").WithCause(err))
		return
	}

	if user == nil {
		log.Error("user not found")
		api.ErrResponse(c, "user not found", errors.NotFound(error_code.ErrForbidden, "login failed"))
		return
	}

	if !tools.CheckPassword(user.Password, req.Password) {
		log.Error("password not match")
		api.ErrResponse(c, "password not match", errors.Forbidden(error_code.ErrForbidden, "login failed"))
		return
	}

	err = repo.CreateLoginLog(repo.GetDb(c), domain.BTMLoginLog{
		UserID:   user.Id,
		UserName: req.Username,
		IP:       c.ClientIP(),
		Browser:  c.Request.UserAgent(),
	})
	if err != nil {
		log.Error("CreateLoginLog error", zap.Any("err", err))
		api.ErrResponse(c, "CreateLoginLog error", errors.InternalServer(error_code.ErrDiError, "CreateLoginLog").WithCause(err))
		return
	}

	token, err := tools.GenerateJWT(domain.UserJwt{
		Account: req.Username,
		Role:    user.Roles,
		Id:      int64(user.Id),
	}, configs.C.JWT.Secret)
	if err != nil {
		log.Error("GenerateJWT", zap.Any("err", err))
		api.ErrResponse(c, "GenerateJWT", errors.InternalServer(error_code.ErrJWT, "GenerateJWT").WithCause(err))
		return
	}

	// 寫入現有有效 token
	err = repo.CreateOrUpdateLastLoginToken(repo.GetDb(c), user.Id, token)
	if err != nil {
		log.Error("CreateOrUpdateLastLoginToken", zap.Any("err", err))
		api.ErrResponse(c, "CreateOrUpdateLastLoginToken", errors.InternalServer(error_code.ErrDiError, "CreateOrUpdateLastLoginToken").WithCause(err))
		return
	}

	api.OKResponse(c, LoginBTMAdminRespItem{
		Token: token,
	})
}
