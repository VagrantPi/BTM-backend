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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type CreateWhitelistReq struct {
	CustomerID uuid.UUID `json:"customer_id" binding:"required"`
	CryptoCode string    `json:"crypto_code" binding:"required"`
	Address    string    `json:"address" binding:"required"`
}

func CreateWhitelist(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "CreateWhitelist")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := CreateWhitelistReq{}
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

	customer, err := repo.GetCustomerById(repo.GetDb(c), req.CustomerID)
	if err != nil {
		log.Error("repo.GetCustomerById", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomerById", errors.NotFound(error_code.ErrDBError, "repo.GetCustomerById").WithCause(err))
		return
	}

	if customer == nil {
		log.Error("customer not found")
		api.ErrResponse(c, "customer not found", errors.NotFound(error_code.ErrDBError, "customer not found"))
		return
	}

	whitelist := &domain.BTMWhitelist{
		CustomerID: req.CustomerID,
		CryptoCode: req.CryptoCode,
		Address:    req.Address,
	}

	tx, err := repo.NewTransactionBegin(c)
	if err != nil {
		log.Error("repo.NewTransactionBegin()", zap.Any("err", err))
		api.ErrResponse(c, "repo.NewTransactionBegin()", errors.InternalServer(error_code.ErrDiError, "repo.NewTransactionBegin()").WithCause(err))
		return
	}
	defer repo.TransactionCommit(tx)

	// create whitelist
	err = repo.CreateWhitelist(tx, whitelist)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		var postgresErr *pgconn.PgError
		if errors.As(err, &postgresErr) && postgresErr.Code == "23505" {
			api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.BadRequest(error_code.ErrWhitelistDuplicate, "duplicated whitelist").WithCause(err))
			return
		}

		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	// add change log
	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	createJsonData, err := json.Marshal(whitelist)
	if err != nil {
		log.Error("json.Marshal(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(whitelist)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(whitelist)").WithCause(err))
		return
	}

	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: operationUserInfo.Id,
		TableName:       domain.BTMChangeLogTableNameBTMWhitelist,
		OperationType:   domain.BTMChangeLogOperationTypeCreate,
		BeforeValue:     nil,
		AfterValue:      createJsonData,
	})
	if err != nil {
		log.Error("CreateBTMChangeLog err", zap.Any("err", err))
		api.ErrResponse(c, "CreateBTMChangeLog", errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: whitelist,
	})
}
