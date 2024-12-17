package controller

import (
	"BTM-backend/internal/controller/config"
	"BTM-backend/internal/controller/customer"
	"BTM-backend/internal/controller/user"
	"BTM-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(apiGroup *gin.RouterGroup) {
	openGroup := apiGroup.Group("/user")
	openGroup.POST("/login", user.LoginBTMAdmin)

	group := apiGroup.Group("/user", middleware.Auth())
	group.GET("/info", user.GetBTMUserInfo)
	group.POST("/logout", user.LogoutBTMAdmin)
	group.GET("/role/routes", user.GetBTMUserRoleRoutes)
	group.GET("/role/roles", user.GetBTMUserRoles)
}

func CustomerRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/customer", middleware.Auth())
	group.GET("/list", customer.GetCustomersList)
	group.GET("/search", customer.SearchCustomers)
	group.GET("/whitelist", customer.GetWhitelist)
	group.GET("/whitelist/search", customer.SearchWhitelist)
	group.POST("/whitelist", customer.CreateWhitelist)
	group.DELETE("/whitelist", customer.DeleteWhitelist)
}

func UserConfigRouter(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/config", middleware.Auth(), config.GetConfig)
}
