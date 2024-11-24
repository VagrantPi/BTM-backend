package controller

import (
	"BTM-backend/internal/controller/user"

	"github.com/gin-gonic/gin"
)

func UserRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/user")
	group.GET("/get", user.GetUser)
}
