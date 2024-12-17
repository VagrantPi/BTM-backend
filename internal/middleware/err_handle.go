package middleware

import (
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type ErrorJsonResp struct {
	Msg    string            `json:"msg"`
	Reason string            `json:"reason"`
	Code   string            `json:"code"`
	Data   map[string]string `json:"data"`
}

func ErrHandler(c *gin.Context) {
	log := logger.Zap().WithClassFunction("system", "ErrHandler")
	defer func() {
		_ = log.Sync()
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Error("gin ErrHandler panic recover:", zap.Error(err.(*errors.Error)))
			obj := ErrorJsonResp{
				Msg:    "internal server error",
				Code:   error_code.ErrInternalPanicError,
				Reason: "",
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, obj)
			return
		}
	}()

	c.Next()

	if c.Errors == nil || c.Errors.Last() == nil {
		return
	}

	err := errors.FromError(c.Errors.Last().Unwrap())
	// 錯誤要印出來給 cloud logging
	logger.Zap().Error("ErrHandler", zap.Any("err", err))

	status := int(err.GetCode())
	// default resp
	obj := ErrorJsonResp{
		Msg:    fmt.Sprintf("無法正常載入，請稍後再試或聯繫客服處理 [%v]", err.GetReason()),
		Code:   err.GetReason(),
		Reason: "請聯繫客服",
	}

	switch status {
	case http.StatusBadRequest:
		obj.Data = err.GetMetadata()
		obj.Msg = err.GetMessage()
		obj.Reason = err.GetReason()
	case http.StatusUnauthorized:
		obj.Msg = "Unauthorized"
	case http.StatusForbidden:
		obj.Msg = "Forbidden"
	}

	c.AbortWithStatusJSON(status, obj)
}
