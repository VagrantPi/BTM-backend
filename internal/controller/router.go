package controller

import (
	"BTM-backend/internal/controller/customer"
	"BTM-backend/internal/controller/user"

	"github.com/gin-gonic/gin"
)

func UserRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/user")
	group.POST("/login", user.LoginBTMAdmin)
	group.GET("/info", user.GetBTMUserInfo)
}

func CustomerRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/customer")
	group.GET("/list", customer.GetCustomersList)
}
