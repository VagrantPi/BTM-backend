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
	Msg    string            `json:"msg"`
	Reason string            `json:"reason"`
	Data   map[string]string `json:"data"`
}

func OKResponse(c *gin.Context, msg string, data map[string]string) {
	if msg == "" {
		msg = "success"
	}
	c.JSON(
		200, ResponseOK{
			Msg:    msg,
			Reason: "",
			Data:   data,
		},
	)
	c.Done()
}

func ErrResponse(c *gin.Context, logInfo string, err error) {
	errUnwrap := errors.FromError(err)
	log.Errorf("%v, err: %v", logInfo, err)
	c.JSON(
		int(errUnwrap.Code), ResponseOK{
			Msg:    errUnwrap.Message,
			Reason: errUnwrap.Reason,
			Data:   errUnwrap.Metadata,
		},
	)
	c.Done()
}
