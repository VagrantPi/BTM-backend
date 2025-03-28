package sumsub

import (
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	sumsubApi "BTM-backend/third_party/sumsub"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetUsersReq struct {
	ApplicantId string `uri:"applicant_id"`
}

func GetApplicantReviewHistory(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetApplicantReviewHistory")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := GetUsersReq{}
	err := c.ShouldBindUri(&req)
	if err != nil {
		log.Error("c.ShouldBindUri(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindUri(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindUri(&req)").WithCause(err))
		return
	}
	fmt.Println("req", req)

	resWithCache, ok := domain.GetTTLMap[domain.UserReviewHistory](&domain.TTLUserHistoryMap, req.ApplicantId)
	if !ok {
		res, err := sumsubApi.GetApplicantReviewHistory(req.ApplicantId)
		if err != nil {
			log.Error("sumsubApi.GetApplicantReviewHistory()", zap.Any("err", err))
			api.ErrResponse(c, "sumsubApi.GetApplicantReviewHistory()", errors.InternalServer(error_code.ErrSumsubApiError, "sumsubApi.GetApplicantReviewHistory()").WithCause(err))
			return
		}

		expire := time.Now().Add(30 * time.Second).UnixNano()
		// cache
		roleWithTTL := domain.TTLMap[domain.UserReviewHistory]{
			Cache: domain.UserReviewHistory{
				CacheHistory: res,
				Expiration:   expire,
			},
			Expire: expire,
		}

		// 使用正確的結構調用 SetTTLMap
		domain.SetTTLMap[domain.UserReviewHistory](&domain.TTLUserHistoryMap, req.ApplicantId, roleWithTTL.Cache, roleWithTTL.Expire)
		resWithCache = &roleWithTTL.Cache
	}

	if resWithCache == nil {
		log.Error("resWithCache is nil")
		api.ErrResponse(c, "resWithCache is nil", errors.InternalServer(error_code.ErrSumsubApiError, "resWithCache is nil"))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: resWithCache.CacheHistory,
	})
	c.Done()
}
