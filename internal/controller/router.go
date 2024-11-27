package controller

import (
	"BTM-backend/internal/controller/user"

	"github.com/gin-gonic/gin"
)

func UserRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/user")
	group.POST("/login", user.LoginBTMAdmin)
	group.GET("/info", user.GetBTMUserInfo)
}
