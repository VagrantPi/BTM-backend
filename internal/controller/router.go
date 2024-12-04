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
	group.GET("/search", customer.SearchCustomers)
	group.GET("/whitelist", customer.GetWhitelist)
	group.GET("/whitelist/search", customer.SearchWhitelist)
	group.POST("/whitelist", customer.CreateWhitelist)
	group.DELETE("/whitelist", customer.DeleteWhitelist)
}
