package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GetCustomerNotesUri struct {
	CustomerID string `uri:"customer_id" binding:"required"`
}

type GetCustomerNotesReq struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

type GetCustomerNotesData struct {
	Total int64                    `json:"total"`
	Items []domain.BTMCustomerNote `json:"items"`
}

func GetCustomerNotes(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetCustomerNotes")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := GetCustomerNotesUri{}
	err := c.ShouldBindUri(&req)
	if err != nil {
		log.Error("c.ShouldBindUri(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindUri(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindUri(&req)").WithCause(err))
		return
	}
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		log.Error("uuid.Parse(req.CustomerID)", zap.Any("err", err))
		api.ErrResponse(c, "uuid.Parse(req.CustomerID)", errors.BadRequest(error_code.ErrInvalidRequest, "uuid.Parse(req.CustomerID)").WithCause(err))
		return
	}

	reqQuery := GetCustomerNotesReq{}
	err = c.BindQuery(&reqQuery)
	if err != nil {
		log.Error("c.BindQuery(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery(&req)").WithCause(err))
		return
	}

	if reqQuery.Limit <= 0 {
		reqQuery.Limit = 10
	}
	if reqQuery.Page <= 0 {
		reqQuery.Page = 1
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	notes, total, err := repo.GetCustomerNotes(repo.GetDb(c), customerID, reqQuery.Limit, reqQuery.Page)
	if err != nil {
		log.Error("repo.GetCustomerNotes", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomerNotes", errors.NotFound(error_code.ErrDBError, "repo.GetCustomerNotes").WithCause(err))
		return
	}

	api.OKResponse(c, GetCustomerNotesData{
		Total: total,
		Items: notes,
	})
}
