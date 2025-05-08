package customer

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

type GetCustomerIdNumberRes struct {
	IdNumber string `json:"id_number"`
}

func GetCustomerIdNumber(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetCustomerIdNumber")
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
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	sumsubInfo, err := repo.GetBTMSumsub(repo.GetDb(c), customerID.String())
	if err != nil {
		log.Error("repo.GetBTMSumsub", zap.Any("customerID", customerID), zap.Any("err", err))
		api.ErrResponse(c, "repo.GetBTMSumsub", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "repo.GetBTMSumsub").WithCause(err))
		return
	}

	// 檢查客製化限額是否存在，不存在則建立
	// 原本設計上未 KYC 不會有客製化限額，所以在撈出時會是 limit: 0，但前端不應顯示 suspend，而應該要走原本 none-kyc 流程
	// 因此查詢時都去做初始化的動作
	_, err = repo.GetRiskControlCustomerLimitSetting(repo.GetDb(c), customerID)
	if err != nil {
		log.Error("repo.GetRiskControlCustomerLimitSetting", zap.Any("customerID", customerID), zap.Any("err", err))
		api.ErrResponse(c, "repo.GetRiskControlCustomerLimitSetting", errors.InternalServer(error_code.ErrDBError, "repo.GetRiskControlCustomerLimitSetting").WithCause(err))
		return
	}

	// 如果沒有存 email 則清除 DB 快取
	if sumsubInfo != nil && (sumsubInfo.EmailHash == "" || sumsubInfo.InspectionId == "") {
		err := repo.DeleteBTMSumsub(repo.GetDb(c), customerID.String())
		if err != nil {
			log.Error("repo.DeleteBTMSumsub", zap.Any("customerID", customerID), zap.Any("err", err))
			api.ErrResponse(c, "repo.DeleteBTMSumsub", errors.InternalServer(error_code.ErrDBError, "repo.DeleteBTMSumsub").WithCause(err))
			return
		}
		sumsubInfo = nil
	}

	if sumsubInfo == nil {
		// fetch sumsub
		idNumber, err := sumsub.FetchDataAdapter(c.Request.Context(), log, repo, customerID.String())
		if err != nil || idNumber == "" {
			log.Error("sumsub.FetchDataAdapter", zap.Any("customerID", customerID), zap.Any("err", err))
			api.ErrResponse(c, "sumsub.FetchDataAdapter", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "sumsub.FetchDataAdapter").WithCause(err))
			return
		}

		api.OKResponse(c, GetCustomerIdNumberRes{
			IdNumber: idNumber,
		})
		return
	}

	exist, fetchExpireDate, err := repo.IsBTMCIBExist(repo.GetDb(c), sumsubInfo.IdNumber)
	if err != nil {
		log.Error("repo.IsBTMCIBExist", zap.Any("customerID", customerID), zap.Any("err", err))
		api.ErrResponse(c, "repo.IsBTMCIBExist", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "repo.IsBTMCIBExist").WithCause(err))
		return
	}
	if exist && !sumsubInfo.BanExpireDate.Valid {
		// 前幾次沒命中，現在命中
		err := repo.UpdateBTMSumsubBanExpireDate(repo.GetDb(c), customerID.String(), fetchExpireDate)
		if err != nil {
			log.Error("repo.UpdateBTMSumsubBanExpireDate", zap.Any("customerID", customerID), zap.Any("err", err))
			api.ErrResponse(c, "repo.UpdateBTMSumsubBanExpireDate", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "repo.UpdateBTMSumsubBanExpireDate").WithCause(err))
			return
		}
	}

	api.OKResponse(c, GetCustomerIdNumberRes{
		IdNumber: sumsubInfo.IdNumber,
	})
}
