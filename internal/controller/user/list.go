package user

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetUsersReq struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

type GetUsersRep struct {
	Total int                       `json:"total"`
	Items []domain.BTMUserWithRoles `json:"items"`
}

func GetUsers(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetUsers")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := GetUsersReq{}
	err := c.BindQuery(&req)
	if err != nil {
		log.Error("c.BindQuery(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery(&req)").WithCause(err))
		return
	}

	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	users, total, err := repo.GetUsers(repo.GetDb(c), req.Limit, req.Page)
	if err != nil {
		log.Error("repo.GetUsers()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetUsers()", errors.InternalServer(error_code.ErrInternalError, "repo.GetUsers()").WithCause(err))
		return
	}

	api.OKResponse(c, GetUsersRep{
		Total: int(total),
		Items: users,
	})
	c.Done()
}
