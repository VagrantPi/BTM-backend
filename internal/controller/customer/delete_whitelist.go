package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type DeleteWhitelistReq struct {
	ID int64 `json:"id" binding:"required"`
}

func DeleteWhitelist(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "DeleteWhitelist")
	defer func() {
		_ = log.Sync()
	}()

	req := DeleteWhitelistReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Error("c.ShouldBindJSON(&req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindJSON(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindJSON(&req)").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	tx, err := repo.NewTransactionBegin(c)
	if err != nil {
		log.Error("repo.NewTransactionBegin()", zap.Any("err", err))
		api.ErrResponse(c, "repo.NewTransactionBegin()", errors.InternalServer(error_code.ErrDiError, "repo.NewTransactionBegin()").WithCause(err))
		return
	}
	defer repo.TransactionCommit(tx)

	// fetch whitelist before delete
	beforeDeleteWhitelist, err := repo.GetWhiteListById(tx, req.ID)
	if err != nil {
		log.Error("repo.GetWhiteListById", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetWhiteListById", errors.InternalServer(error_code.ErrDBError, "repo.GetWhitelistById").WithCause(err))
		return
	}

	// delete whitelist
	err = repo.DeleteWhitelist(tx, req.ID)
	if err != nil {
		log.Error("repo.DeleteWhitelist", zap.Any("err", err))
		api.ErrResponse(c, "repo.DeleteWhitelist", errors.InternalServer(error_code.ErrDBError, "repo.DeleteWhitelist").WithCause(err))
		return
	}

	// create change log
	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	beforeDeleteWhitelistJsonData, err := json.Marshal(beforeDeleteWhitelist)
	if err != nil {
		log.Error("json.Marshal(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(whitelist)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(whitelist)").WithCause(err))
		return
	}

	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: operationUserInfo.Id,
		TableName:       domain.BTMChangeLogTableNameBTMWhitelist,
		OperationType:   domain.BTMChangeLogOperationTypeDelete,
		CustomerId:      beforeDeleteWhitelist.CustomerID,
		BeforeValue:     beforeDeleteWhitelistJsonData,
		AfterValue:      nil,
	})
	if err != nil {
		log.Error("CreateBTMChangeLog err", zap.Any("err", err))
		api.ErrResponse(c, "CreateBTMChangeLog", errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
	})
}
