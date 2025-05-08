package tool

import (
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/third_party/sms"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetSmsHistoryReq struct {
	DateAt time.Time `form:"date_at"`
	Limit  int       `form:"limit"`
}

type GetSmsHistoryRep struct {
	Items []domain.TwilioHistorySms `json:"items"`
}

func GetSmsHistory(c *gin.Context) {

	log := logger.Zap().WithClassFunction("api", "GetSmsHistory")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	var req GetSmsHistoryReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
		return
	}
	if req.DateAt.IsZero() {
		req.DateAt = time.Now().Add(-1 * time.Hour)
	}

	if req.Limit == 0 || req.Limit > 300 {
		req.Limit = 30
	}

	smsHistory, err := sms.GetSmsHistory(req.DateAt, req.Limit)
	if err != nil {
		log.Error("GetSmsHistory error", zap.Any("err", err))
		api.ErrResponse(c, "GetSmsHistory error", errors.InternalServer(error_code.ErrThirdPartyHttpCall, "GetSmsHistory error").WithCause(err))
		return
	}

	api.OKResponse(c, GetSmsHistoryRep{
		Items: smsHistory,
	})
}
