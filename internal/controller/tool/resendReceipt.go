package tool

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/third_party/sms"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type ResendReceiptReq struct {
	SessionId string `uri:"session_id"`
}

func ResendReceipt(c *gin.Context) {

	log := logger.Zap().WithClassFunction("api", "ResendReceipt")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	// 撈出未被軟刪除的地址綁定
	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	var req ResendReceiptReq
	err = c.ShouldBindUri(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
		return
	}

	fetch, err := repo.GetCashInTxBySessionId(repo.GetDb(c), req.SessionId)
	if err != nil {
		log.Error("repo.GetCashInTxBySessionId()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCashInTxBySessionId()", errors.InternalServer(error_code.ErrDBError, "repo.GetCashInTxBySessionId()").WithCause(err))
		return
	}

	if fetch == nil {
		log.Error("repo.GetCashInTxBySessionId()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCashInTxBySessionId()", errors.InternalServer(error_code.ErrDBError, "not found"))
		return
	}

	btcDecimal := decimal.NewFromInt(10).Pow(decimal.NewFromInt(8))
	ethDecimal := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))
	cryptoAtoms := fetch.CryptoAtoms.Div(btcDecimal)
	switch fetch.CryptoCode {
	case "ETH":
		cryptoAtoms = fetch.CryptoAtoms.Div(ethDecimal)
	}

	err = sms.SendSms(fetch.Phone, req.SessionId, fetch.SendTime.Add(time.Hour*8).Format("2006-01-02 15:04:05"), fetch.Fiat.String(), cryptoAtoms.String(), fetch.CryptoCode, fetch.RawTickerPrice.Add(fetch.RawTickerPrice.Mul(fetch.CommissionPercentage)).StringFixed(2), fetch.TxHash)
	if err != nil {
		log.Error("sms.SendSms()", zap.Any("err", err))
		api.ErrResponse(c, "sms.SendSms()", errors.InternalServer(error_code.ErrThirdPartyHttpCall, "sms.SendSms()").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}
