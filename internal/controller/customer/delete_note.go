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
	"go.uber.org/zap"
)

type DeleteCustomerNoteUri struct {
	CustomerId string `uri:"customer_id" binding:"required"`
	NoteId     uint   `uri:"note_id" binding:"required"`
}

func DeleteCustomerNote(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "DeleteCustomerNote")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := DeleteCustomerNoteUri{}
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

	fetchNote, err := repo.GetCustomerNote(tx, req.NoteId)
	if err != nil {
		log.Error("repo.GetCustomerNote(tx, req.NoteId)", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomerNote(tx, req.NoteId)", errors.InternalServer(error_code.ErrDBError, "repo.GetCustomerNote(tx, req.NoteId)").WithCause(err))
		return
	}
	beforeJsonData, err := json.Marshal(fetchNote)
	if err != nil {
		log.Error("json.Marshal(fetchNote)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(fetchNote)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(fetchNote)").WithCause(err))
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
		OperationType:   domain.BTMChangeLogOperationTypeDelete,
		CustomerId:      &customerId,
		BeforeValue:     beforeJsonData,
		AfterValue:      nil,
	})
	if err != nil {
		log.Error("CreateBTMChangeLog err", zap.Any("err", err))
		api.ErrResponse(c, "CreateBTMChangeLog", errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err))
		return
	}

	// delete note
	err = repo.DeleteCustomerNote(tx, req.NoteId)
	if err != nil {
		log.Error("DeleteCustomerNote err", zap.Any("err", err))
		api.ErrResponse(c, "DeleteCustomerNote", errors.InternalServer(error_code.ErrDBError, "DeleteCustomerNote").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}
