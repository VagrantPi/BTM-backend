package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/third_party/sumsub"
	"database/sql"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var fetchSumsubLock *sync.Mutex

func init() {
	fetchSumsubLock = &sync.Mutex{}
}

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

	// 如果沒有存 email 則清除 DB 快取
	if sumsubInfo != nil && sumsubInfo.Email == "" {
		err := repo.DeleteBTMSumsub(repo.GetDb(c), customerID.String())
		if err != nil {
			log.Error("repo.DeleteBTMSumsub", zap.Any("customerID", customerID), zap.Any("err", err))
			api.ErrResponse(c, "repo.DeleteBTMSumsub", errors.InternalServer(error_code.ErrDBError, "repo.DeleteBTMSumsub").WithCause(err))
			return
		}
		sumsubInfo = nil
	}

	if sumsubInfo == nil {
		// cache lock
		fetchSumsubLock.Lock()
		defer fetchSumsubLock.Unlock()

		// fetch sumsub
		data, err := sumsub.GetApplicantInfo(customerID.String())
		if err != nil {
			log.Error("sumsub.GetApplicantInfo", zap.Any("customerID", customerID), zap.Any("err", err))
			api.ErrResponse(c, "sumsub.GetApplicantInfo", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "sumsub.GetApplicantInfo").WithCause(err))
			return
		}

		if data.Review.ReviewStatus != "completed" {
			log.Error("sumsub.GetApplicantInfo ReviewStatus not completed", zap.Any("customerID", customerID), zap.Any("data", data))
			api.ErrResponse(c, "sumsub.GetApplicantInfo", errors.InternalServer(error_code.ErrBTMSumsubIdNumberNotFound, "review status not completed"))
			return
		}

		// 抓取 sumsub id number
		idNumber := ""
		if len(data.Info.IdDocs) > 0 {
			for _, v := range data.Info.IdDocs {
				if v.IdDocType == "ID_CARD" {
					idNumber = v.Number
					break
				}
			}
		}

		if idNumber == "" {
			log.Error("idNumber is empty", zap.Any("customerID", customerID), zap.Any("idNumber", idNumber))
			api.ErrResponse(c, "idNumber is empty", errors.InternalServer(error_code.ErrBTMSumsubIdNumberNotFound, "idNumber is empty"))
			return
		}

		insertData := domain.BTMSumsub{
			ApplicantId:   data.Id,
			CustomerId:    customerID,
			Info:          data,
			IdNumber:      idNumber,
			BanExpireDate: sql.NullInt64{},
			Email:         data.Email,
			Phone:         data.Phone,
		}
		exist, fetchExpireDate, err := repo.IsBTMCIBExist(repo.GetDb(c), idNumber)
		if err != nil {
			log.Error("repo.IsBTMCIBExist", zap.Any("customerID", customerID), zap.Any("err", err))
			api.ErrResponse(c, "repo.IsBTMCIBExist", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "repo.IsBTMCIBExist").WithCause(err))
			return
		}
		if exist {
			insertData.BanExpireDate = sql.NullInt64{
				Int64: fetchExpireDate,
				Valid: true,
			}
		}

		// store db
		err = repo.CreateBTMSumsub(repo.GetDb(c), insertData)
		if err != nil {
			log.Error("repo.CreateBTMSumsub", zap.Any("customerID", customerID), zap.Any("err", err))
			api.ErrResponse(c, "repo.CreateBTMSumsub", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "repo.CreateBTMSumsub").WithCause(err))
			return
		}

		c.JSON(200, api.DefaultRep{
			Code: 20000,
			Data: GetCustomerIdNumberRes{
				IdNumber: idNumber,
			},
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

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: GetCustomerIdNumberRes{
			IdNumber: sumsubInfo.IdNumber,
		},
	})
}
