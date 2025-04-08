package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/third_party/sumsub"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AddSumsubTagRes struct{}

func AddSumsubTag(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "AddSumsubTag")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	// externalUserId(customers.id)
	customerIDStr := c.Query("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		log.Error("uuid.Parse(customerIDStr)", zap.Any("customerID", customerID), zap.Any("err", err))
		api.ErrResponse(c, "uuid.Parse(customerIDStr)", errors.BadRequest(error_code.ErrInvalidRequest, "Invalid UUID format").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("customerID", customerID), zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	sumsubInfo, err := repo.GetBTMSumsub(repo.GetDb(c), customerID.String())
	if err != nil {
		log.Error("repo.GetBTMSumsub", zap.Any("customerID", customerID), zap.Any("err", err))
		api.ErrResponse(c, "repo.GetBTMSumsub", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "repo.GetBTMSumsub").WithCause(err))
		return
	}
	if sumsubInfo == nil {
		log.Error("repo.GetBTMSumsub not found", zap.Any("customerID", customerID))
		api.ErrResponse(c, "repo.GetBTMSumsub not found", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "repo.GetBTMSumsub not found"))
		return
	}

	// add sumsub tag
	err = sumsub.AddAndOverwriteApplicantTags(sumsubInfo.ApplicantId, []string{domain.SumsubTagCib.String()})
	if err != nil {
		log.Error("sumsub.AddAndOverwriteApplicantTags", zap.Any("customerID", customerID), zap.Any("err", err))
		api.ErrResponse(c, "sumsub.AddAndOverwriteApplicantTags", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "sumsub.AddAndOverwriteApplicantTags").WithCause(err))
		return
	}

	api.OKResponse(c, AddSumsubTagRes{})
}
