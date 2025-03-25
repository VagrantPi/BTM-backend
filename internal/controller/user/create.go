package user

import (
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type CreateOneReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     int64  `json:"role" binding:"required"`
}

func CreateOne(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "CreateOne")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	var req CreateOneReq
	err := c.BindJSON(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
		return
	}

	passwordRegexp, err := regexp.Compile("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])[a-zA-Z0-9!@#$%^&*]{10,}$")
	if err != nil {
		log.Error("error parsing regexp", zap.Any("err", err))
		api.ErrResponse(c, "error parsing regexp", errors.InternalServer(error_code.ErrInternalError, "error parsing regexp").WithCause(err))
		return
	}
	if !passwordRegexp.MatchString(req.Password) {
		log.Error("password not match regex", zap.Any("password", req.Password))
		api.ErrResponse(c, "password not match regex", errors.BadRequest(error_code.ErrInvalidRequest, "password not match regex"))
		return
	}

}
