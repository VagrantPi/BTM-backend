package sumsub

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/third_party/sumsub"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type SumsubApplicantWebhook struct {
	ApplicantId    string `json:"applicantId"`
	InspectionId   string `json:"inspectionId"`
	CorrelationId  string `json:"correlationId"`
	LevelName      string `json:"levelName"`
	ExternalUserId string `json:"externalUserId"`
	Type           string `json:"type"`
	SandboxMode    bool   `json:"sandboxMode"`
	ReviewStatus   string `json:"reviewStatus"`
	CreatedAtMs    string `json:"createdAtMs"`
	ApplicantType  string `json:"applicantType"`
	ReviewResult   struct {
		ReviewAnswer     string   `json:"reviewAnswer,omitempty"`
		RejectLabels     []string `json:"rejectLabels,omitempty"`
		ReviewRejectType string   `json:"reviewRejectType,omitempty"`
	} `json:"reviewResult"`
	VideoIdentReviewStatus    string `json:"videoIdentReviewStatus"`
	ApplicantActionId         string `json:"applicantActionId"`
	ExternalApplicantActionId string `json:"externalApplicantActionId"`
	ClientId                  string `json:"clientId"`
}

type OutputSumsubApplicantWebhook struct {
	ApplicantId string `json:"applicant_id"`
}

var apiLock *sync.Mutex
var apiLockInitialFlag sync.Once

func SumsubWebhook(c *gin.Context) {
	fmt.Println("SumsubWebhook!!")
	// init cache lock
	apiLockInitialFlag.Do(func() {
		apiLock = new(sync.Mutex)
	})

	log := logger.Zap().WithClassFunction("api", "SumsubWebhook")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	// 限制同一時間只能處理一個
	apiLock.Lock()
	defer apiLock.Unlock()

	var req SumsubApplicantWebhook
	err := c.BindJSON(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	// 審核後或標籤更新
	if (req.Type == domain.SumsubWebhookTypeApplicantReviewed.String() || req.Type == domain.SumsubWebhookTypeApplicantTagsChanged.String()) &&
		req.ReviewStatus == domain.SumsubApplicantStatusCompleted.String() {

		// fetch sumsub
		_, err := sumsub.FetchDataAdapter(c, log, repo, req.ExternalUserId)
		if err != nil {
			log.Error("sumsub.FetchDataAdapter", zap.Any("customerID", req.ExternalUserId), zap.Any("err", err))
			api.ErrResponse(c, "sumsub.FetchDataAdapter", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "sumsub.FetchDataAdapter").WithCause(err))
			return
		}
	}

	c.JSON(200, OutputSumsubApplicantWebhook{
		ApplicantId: req.ApplicantId,
	})
	c.Done()
}
