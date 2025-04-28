package user

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetUsersLiteRepItem struct {
	Id      uint   `json:"id"`
	Account string `json:"account"`
}

type GetUsersLiteRep struct {
	Items []GetUsersLiteRepItem `json:"items"`
}

func GetUsersLite(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetUsersLite")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	users, _, err := repo.GetUsers(repo.GetDb(c), 1000, 1)
	if err != nil {
		log.Error("repo.GetUsers()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetUsers()", errors.InternalServer(error_code.ErrInternalError, "repo.GetUsers()").WithCause(err))
		return
	}

	var items []GetUsersLiteRepItem
	for _, user := range users {
		items = append(items, GetUsersLiteRepItem{
			Id:      user.Id,
			Account: user.Account,
		})
	}

	api.OKResponse(c, GetUsersLiteRep{
		Items: items,
	})
	c.Done()
}
