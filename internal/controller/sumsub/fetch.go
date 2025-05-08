package sumsub

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/third_party/sumsub"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FetchDataAdapterReq struct {
	CustomerId string `uri:"customer_id"`
}

func FetchDataAdapter(c *gin.Context) {

	log := logger.Zap().WithClassFunction("api", "FetchDataAdapter")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	var req FetchDataAdapterReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
		return
	}

	customerId, err := uuid.Parse(req.CustomerId)
	if err != nil {
		log.Error("uuid.Parse(req.CustomerId)", zap.Any("err", err))
		api.ErrResponse(c, "uuid.Parse(req.CustomerId)", errors.BadRequest(error_code.ErrInvalidRequest, "uuid.Parse(req.CustomerId)").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	// fetch sumsub
	_, err = sumsub.FetchDataAdapter(c, log, repo, customerId.String())
	if err != nil {
		log.Error("sumsub.FetchDataAdapter", zap.Any("customerID", customerId), zap.Any("err", err))
		api.ErrResponse(c, "sumsub.FetchDataAdapter", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "sumsub.FetchDataAdapter").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
	c.Done()
}
