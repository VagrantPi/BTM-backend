package user

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetUser(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetUser")
	defer func() {
		_ = log.Sync()
	}()

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		_ = c.Error(err)
		return
	}
	comstomer, err := repo.GetCustomerByPhone("")
	if err != nil {
		log.Error("GetCustomerByPhone", zap.Any("err", err))
		_ = c.Error(err)
		return
	}

	fmt.Println("comstomer", *comstomer)

	api.OKResponse(c, "", map[string]string{})
}
