package cib

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetCibsListReq struct {
	Id    string `form:"id"`
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
}

type GetCibsListRepItem struct {
	Total int64           `json:"total"`
	Items []domain.BTMCIB `json:"items"`
}

func GetCibsList(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetCibsList")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := GetCibsListReq{}
	err := c.BindQuery(&req)
	if err != nil {
		log.Error("c.BindQuery(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery(&req)").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	cibs, total, err := repo.GetBTMCIBs(repo.GetDb(c), req.Id, req.Limit, req.Page)
	if err != nil {
		log.Error("repo.GetCustomers()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomers()", errors.NotFound(error_code.ErrDBError, "repo.GetCustomers()").WithCause(err))
		return
	}

	api.OKResponse(c, GetCibsListRepItem{
		Total: total,
		Items: cibs,
	})
}
