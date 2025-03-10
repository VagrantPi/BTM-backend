package riskControl

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

func GetRiskControlRoles(c *gin.Context) {
	repo, err := di.NewRepo()
	if err != nil {
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	roles, err := repo.GetRiskControlRoles()
	if err != nil {
		api.ErrResponse(c, "repo.GetRiskControlRoles()", errors.InternalServer(error_code.ErrDBError, "repo.GetRiskControlRoles()").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: roles,
	})
	c.Done()
}
