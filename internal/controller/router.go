package controller

import (
	"BTM-backend/internal/controller/cib"
	"BTM-backend/internal/controller/config"
	"BTM-backend/internal/controller/customer"
	"BTM-backend/internal/controller/debug"
	"BTM-backend/internal/controller/tx"
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
	group.GET("/search/whitelist_created_at", customer.SearchCustomersByWhitelistCreatedAt)
	group.GET("/search/address/:address", customer.SearchCustomersByAddress)
	group.GET("/whitelist", customer.GetWhitelist)
	group.GET("/whitelist/search", customer.SearchWhitelist)
	group.POST("/whitelist", customer.CreateWhitelist)
	group.DELETE("/whitelist", customer.DeleteWhitelist)
}

func CustomerInternalRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/customer")
	group.GET("/id_number", customer.GetCustomerIdNumber)
}

func UserConfigRouter(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/config", middleware.Auth(), config.GetConfig)
}

func TxRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/tx", middleware.Auth())
	group.GET("/list", tx.GetTxsList)
}

// TODO: 未來增加安全性 middleware
func DebugRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/debug", middleware.Auth())
	group.GET("/logs", debug.GetBTMChangeLogs)
}

func CibRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/cib", middleware.Auth())
	group.GET("/list", cib.GetCibsList)
}
