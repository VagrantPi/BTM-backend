package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

func Ping(c *gin.Context) {
	OKResponse(c, "success", nil)
	c.Done()
}

type ResponseOK struct {
	Code   int               `json:"code"`
	Msg    string            `json:"msg"`
	Reason string            `json:"reason"`
	Data   map[string]string `json:"data"`
}

func OKResponse(c *gin.Context, msg string, data map[string]string) {
	if msg == "" {
		msg = "success"
	}
	c.JSON(200, ResponseOK{
		Code:   20000,
		Msg:    msg,
		Reason: "",
		Data:   data,
	})
	c.Done()
}

func ErrResponse(c *gin.Context, logInfo string, err error) {
	c.Set("custom_error", err)
	errUnwrap := errors.FromError(err)
	log.Errorf("%v, err: %v", logInfo, err)
	c.JSON(int(errUnwrap.Code), ResponseOK{
		Code:   20000,
		Msg:    errUnwrap.Message,
		Reason: errUnwrap.Reason,
		Data:   errUnwrap.Metadata,
	})
	c.Done()
}

type DefaultRep struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
