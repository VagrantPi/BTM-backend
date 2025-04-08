package tool

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

func CompleteAddressBindingLog(c *gin.Context) {

	log := logger.Zap().WithClassFunction("api", "CompleteAddressBindingLog")
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

	whitelists, err := repo.GetWhiteLists(repo.GetDb(c))
	if err != nil {
		log.Error("repo.GetWhiteLists()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetWhiteLists()", errors.InternalServer(error_code.ErrDBError, "repo.GetWhiteLists()").WithCause(err))
		return
	}

	lostLogs := []domain.BTMChangeLog{}
	for _, whitelist := range whitelists {
		exists, err := repo.IsAddressExistsInAfterValue(repo.GetDb(c), whitelist.Address)
		if err != nil {
			log.Error("repo.IsAddressExistsInAfterValue()", zap.Any("err", err))
			api.ErrResponse(c, "repo.IsAddressExistsInAfterValue()", errors.InternalServer(error_code.ErrDBError, "repo.IsAddressExistsInAfterValue()").WithCause(err))
			return
		}
		if !exists {
			whitelistJson := &domain.BTMWhitelist{
				CustomerID: whitelist.CustomerID,
				CryptoCode: whitelist.CryptoCode,
				Address:    whitelist.Address,
			}
			createJsonData, err := json.Marshal(whitelistJson)
			if err != nil {
				log.Error("json.Marshal(whitelist)", zap.Any("err", err))
				api.ErrResponse(c, "json.Marshal(whitelist)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(whitelist)").WithCause(err))
				return
			}

			lostLogs = append(lostLogs, domain.BTMChangeLog{
				OperationUserId: 0,
				TableName:       domain.BTMChangeLogTableNameBTMWhitelist,
				OperationType:   domain.BTMChangeLogOperationTypeCreate,
				CustomerId:      &whitelist.CustomerID,
				BeforeValue:     nil,
				AfterValue:      createJsonData,
				CreatedAt:       whitelist.CreatedAt,
			})
		}
	}

	log.Info("補綁定地址紀錄", zap.Int("資料筆數", len(lostLogs)))

	if len(lostLogs) > 0 {
		err = repo.BatchCreateBTMChangeLog(repo.GetDb(c), lostLogs)
		if err != nil {
			log.Error("repo.BatchCreateBTMChangeLog()", zap.Any("err", err))
			api.ErrResponse(c, "repo.BatchCreateBTMChangeLog()", errors.InternalServer(error_code.ErrDBError, "repo.BatchCreateBTMChangeLog()").WithCause(err))
			return
		}
	}

	api.OKResponse(c, nil)
}
