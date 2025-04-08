package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

func Ping(c *gin.Context) {
	OKResponse(c, nil)
	c.Done()
}

func OKResponse(c *gin.Context, data interface{}) {
	c.JSON(200, DefaultRep{
		Code:      20000,
		Msg:       "success",
		Reason:    "",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
	c.Done()
}

func ErrResponse(c *gin.Context, logInfo string, err error) {
	c.Set("custom_error", err)
	errUnwrap := errors.FromError(err)
	log.Errorf("%v, err: %v", logInfo, err)
	c.JSON(int(errUnwrap.Code), DefaultRep{
		Code:      20000,
		Msg:       errUnwrap.Message,
		Reason:    errUnwrap.Reason,
		Data:      errUnwrap.Metadata,
		Timestamp: time.Now().Unix(),
	})
	c.Done()
}

type DefaultRep struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Reason    string      `json:"reason"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}
