package customer

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
		log.Error("uuid.Parse(customerIDStr)", zap.Any("err", err))
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
		log.Error("repo.GetBTMSumsub", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetBTMSumsub", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "repo.GetBTMSumsub").WithCause(err))
		return
	}
	if sumsubInfo == nil {
		// cache lock
		fetchSumsubLock.Lock()
		defer fetchSumsubLock.Unlock()

		// fetch sumsub
		data, err := sumsub.GetApplicantInfo(customerID.String())
		if err != nil {
			log.Error("sumsub.GetApplicantInfo", zap.Any("err", err))
			api.ErrResponse(c, "sumsub.GetApplicantInfo", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "sumsub.GetApplicantInfo").WithCause(err))
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
			log.Error("idNumber is empty")
			api.ErrResponse(c, "idNumber is empty", errors.InternalServer(error_code.ErrBTMSumsubIdNumberNotFound, "idNumber is empty"))
			return
		}

		fmt.Println("idNumber", idNumber)
		fmt.Println("data", data)
		// store db
		err = repo.CreateBTMSumsub(repo.GetDb(c), domain.BTMSumsub{
			CustomerId: customerID,
			Info:       data,
			IdNumber:   idNumber,
		})
		if err != nil {
			log.Error("repo.CreateBTMSumsub", zap.Any("err", err))
			api.ErrResponse(c, "repo.StoreBTMSumsub", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "repo.StoreBTMSumsub").WithCause(err))
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

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: GetCustomerIdNumberRes{
			IdNumber: sumsubInfo.IdNumber,
		},
	})
}
