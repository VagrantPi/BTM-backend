package tool

import (
	"BTM-backend/internal/cronjob"
	"BTM-backend/pkg/api"

	"github.com/gin-gonic/gin"
)

func SyncAllCustomerSumsub(c *gin.Context) {
	go cronjob.SyncNotComplateSumsub()

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: "success",
	})
	c.Done()
}
