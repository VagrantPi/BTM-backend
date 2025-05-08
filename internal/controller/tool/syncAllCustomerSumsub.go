package tool

import (
	"BTM-backend/internal/cronjob"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

type SyncAllCustomerSumsubReq struct {
	Force bool `form:"force"`
}

func SyncAllCustomerSumsub(c *gin.Context) {

	var req SyncAllCustomerSumsubReq
	if err := c.ShouldBindQuery(&req); err != nil {
		api.ErrResponse(c, "BindQuery error", errors.BadRequest(error_code.ErrInvalidRequest, "BindQuery error").WithCause(err))
		return
	}

	// 當 force 為 true 時，幾乎全部用戶都會跑，呼叫時需留意
	go cronjob.SyncNotComplateSumsub(req.Force)

	api.OKResponse(c, nil)
}
