package controller

import (
	"BTM-backend/internal/controller/cib"
	"BTM-backend/internal/controller/config"
	"BTM-backend/internal/controller/customer"
	"BTM-backend/internal/controller/debug"
	"BTM-backend/internal/controller/riskControl"
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
	group.POST("/role", user.CreateRole)
	group.PUT("/role", user.UpdateRole)
}

func CustomerRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/customer", middleware.Auth())
	group.GET("/list", customer.SearchCustomers)
	group.GET("/whitelist", customer.GetWhitelist)
	group.GET("/whitelist/search", customer.SearchWhitelist)
	group.POST("/whitelist", customer.CreateWhitelist)
	group.DELETE("/whitelist", customer.DeleteWhitelist)
}

func UserConfigRouter(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/config", middleware.Auth(), config.GetConfig)
}

func TxRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/tx", middleware.Auth())
	group.GET("/list", tx.GetTxsList)
}

func CibRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/cib", middleware.Auth())
	group.GET("/list", cib.GetCibsList)
	group.POST("/upload", cib.UploadCib)
}

func RiskControlRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/risk_control", middleware.Auth())
	group.GET("/roles", riskControl.GetRiskControlRoles)
	group.GET("/:customer_id/role", riskControl.GetCustomerRiskControlRole)
	group.PATCH("/:customer_id/role", riskControl.UpdateCustomerRiskControlRole)
	group.PATCH("/:customer_id/limit", riskControl.UpdateCustomerRiskControlLimit)
}

func InternalRouter(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/btm/logs", middleware.Auth(), debug.GetBTMChangeLogs)

	group := apiGroup.Group("/btm", middleware.ServerKeyAuth())
	group.GET("/id_number", customer.GetCustomerIdNumber)
	group.POST("/cib", debug.DownlaodCIB)
	group.POST("/add_sumsub_tag", customer.AddSumsubTag)
}
