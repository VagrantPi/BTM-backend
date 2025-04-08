package tool

import (
	"BTM-backend/internal/cronjob"
	"BTM-backend/pkg/api"

	"github.com/gin-gonic/gin"
)

func SyncAllCustomerSumsub(c *gin.Context) {
	go cronjob.SyncNotComplateSumsub()

	api.OKResponse(c, nil)
}
