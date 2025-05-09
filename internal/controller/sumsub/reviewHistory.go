package sumsub

import (
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	sumsubApi "BTM-backend/third_party/sumsub"
	"sort"
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

	api.OKResponse(c, parseHistory(resWithCache.CacheHistory.Items))
}

type GetApplicantReviewHistoryItem struct {
	SubmissionTime string   `json:"submissionTime"`
	ReviewTime     string   `json:"reviewTime"`
	Status         string   `json:"status"`
	Reason         []string `json:"reason"`
}

func parseHistory(items []domain.SumsubHistoryItem) []GetApplicantReviewHistoryItem {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Ts < items[j].Ts
	})

	var submissionTimes []string

	res := []GetApplicantReviewHistoryItem{}

	for i := 0; i < len(items); i++ {
		item := items[i]
		// SubjectName: coin_now.com 為系統自動審核
		if item.Status == "completed" && item.SubjectName != "coin_now.com" {
			status := ""
			if item.ReviewAnswer == "GREEN" && item.Activity != "user:changed:applicantTag" {
				status = "已核准"
			} else if item.ReviewAnswer == "RED" && (item.ReviewRejectType == "FINAL" || item.ReviewRejectType == "RETRY") {
				status = "已拒絕"
			}

			if status != "" {
				reviewTime := item.Ts
				submissionTime := reviewTime
				if len(submissionTimes) > 0 {
					submissionTime = submissionTimes[len(submissionTimes)-1]
					submissionTimes = nil
				}
				res = append(res, GetApplicantReviewHistoryItem{
					SubmissionTime: submissionTime,
					ReviewTime:     reviewTime,
					Status:         status,
					Reason:         getReason(item.ReviewResult.RejectLabels),
				})
			}
		} else {
			submissionTimes = append([]string{item.Ts}, submissionTimes...)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ReviewTime > res[j].ReviewTime
	})

	return res
}

func getReason(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}
	var reasons []string
	for _, tag := range tags {
		kycTag := domain.KYCTag[tag]
		if kycTag == "" {
			kycTag = tag
		}
		reasons = append(reasons, kycTag)
	}
	return reasons
}
