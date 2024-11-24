package user

import (
	"BTM-backend/pkg/api"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	api.OKResponse(c, "success")
}
