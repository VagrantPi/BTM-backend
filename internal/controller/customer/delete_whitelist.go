package customer

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type DeleteWhitelistReq struct {
	ID int64 `json:"id" binding:"required"`
}

type DeleteWhitelistRep struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
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

	err = repo.DeleteWhitelist(req.ID)
	if err != nil {
		log.Error("repo.DeleteWhitelist", zap.Any("err", err))
		api.ErrResponse(c, "repo.DeleteWhitelist", errors.InternalServer(error_code.ErrDBError, "repo.DeleteWhitelist").WithCause(err))
		return
	}

	c.JSON(200, CreateWhitelistRep{
		Code: 20000,
	})
}
