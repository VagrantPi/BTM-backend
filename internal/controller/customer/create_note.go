package customer

import (
	"BTM-backend/configs"
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
	"go.uber.org/zap"
)

type CreateCustomerNoteUri struct {
	CustomerId string `uri:"customer_id" binding:"required"`
}

type CreateCustomerNoteReq struct {
	Note     string `json:"note" binding:"required"`
	NoteType int    `json:"note_type" binding:"required"`
}

func CreateCustomerNote(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "CreateCustomerNote")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := CreateCustomerNoteUri{}
	err := c.ShouldBindUri(&req)
	if err != nil {
		log.Error("c.ShouldBindUri(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindUri(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindUri(&req)").WithCause(err))
		return
	}
	customerId, err := uuid.Parse(req.CustomerId)
	if err != nil {
		log.Error("uuid.Parse(req.CustomerId)", zap.Any("err", err))
		api.ErrResponse(c, "uuid.Parse(req.CustomerId)", errors.BadRequest(error_code.ErrInvalidRequest, "uuid.Parse(req.CustomerId)").WithCause(err))
		return
	}

	reqBody := CreateCustomerNoteReq{}
	err = c.ShouldBindJSON(&reqBody)
	if err != nil {
		log.Error("c.ShouldBindJSON(&req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindJSON(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindJSON(&req)").WithCause(err))
		return
	}

	repo, err := di.NewRepo(configs.C.Mock)
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

	input := domain.BTMCustomerNote{
		CustomerId: customerId,
		Note:       reqBody.Note,
		NoteType:   domain.CustomerNoteType(reqBody.NoteType),
	}
	afterJsonData, err := json.Marshal(input)
	if err != nil {
		log.Error("json.Marshal(input)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(input)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(input)").WithCause(err))
		return
	}

	// add change log
	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: operationUserInfo.Id,
		TableName:       domain.BTMChangeLogTableNameBTMCustomerNotes,
		OperationType:   domain.BTMChangeLogOperationTypeCreate,
		CustomerId:      &customerId,
		BeforeValue:     nil,
		AfterValue:      afterJsonData,
	})
	if err != nil {
		log.Error("CreateBTMChangeLog err", zap.Any("err", err))
		api.ErrResponse(c, "CreateBTMChangeLog", errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err))
		return
	}
	input.OperationUserId = operationUserInfo.Id
	input.OperationUserName = operationUserInfo.Account

	// create note
	err = repo.CreateCustomerNote(tx, input)
	if err != nil {
		log.Error("CreateCustomerNote err", zap.Any("err", err))
		api.ErrResponse(c, "CreateCustomerNote", errors.InternalServer(error_code.ErrDBError, "CreateCustomerNote").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}
