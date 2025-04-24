package controller

import (
	"BTM-backend/internal/controller/cib"
	"BTM-backend/internal/controller/config"
	"BTM-backend/internal/controller/customer"
	"BTM-backend/internal/controller/debug"
	"BTM-backend/internal/controller/riskControl"
	"BTM-backend/internal/controller/sumsub"
	"BTM-backend/internal/controller/tool"
	"BTM-backend/internal/controller/tx"
	"BTM-backend/internal/controller/user"
	"BTM-backend/internal/controller/view"
	"BTM-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(apiGroup *gin.RouterGroup) {
	openGroup := apiGroup.Group("/user")
	openGroup.POST("/login", user.LoginBTMAdmin)

	group := apiGroup.Group("/user", middleware.Auth(), middleware.CheckRole())
	group.GET("/info", user.GetBTMUserInfo)
	group.POST("/logout", user.LogoutBTMAdmin)
	group.POST("/one", user.CreateOne)
	group.PUT("/one", user.UpdateOne)
	group.GET("/list", user.GetUsers)
	group.GET("/list/lite", user.GetUsersLite)

	// user role
	group.GET("/role/routes", user.GetBTMUserRoleRoutes)
	group.GET("/role/roles", user.GetBTMUserRoles)
	group.POST("/role", user.CreateRole)
	group.PUT("/role", user.UpdateRole)

}

func ThirdPartyRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/3rd")

	group.POST("/v1/sumsub/webhook", middleware.SumsubGuardImpl.CheckDigest, sumsub.SumsubWebhook)
	group.GET("/v1/sumsub/:applicant_id/review_history", middleware.Auth(), middleware.CheckRole(), sumsub.GetApplicantReviewHistory)
}

func CustomerRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/customer", middleware.Auth(), middleware.CheckRole())
	group.GET("/list", customer.SearchCustomers)
	group.GET("/list/edd", customer.SearchEddUsers)
	group.GET("/:customer_id/detail", customer.GetBTMUserInfoDetail)
	group.GET("/whitelist", customer.GetWhitelist)
	group.GET("/whitelist/search", customer.SearchWhitelist)
	group.POST("/whitelist", customer.CreateWhitelist)
	group.DELETE("/whitelist", customer.DeleteWhitelist)
	group.GET("/image", customer.GetSumsubImage)
	group.GET("/:customer_id/notes", customer.GetCustomerNotes)
	group.POST("/:customer_id/note", customer.CreateCustomerNote)
}

func UserConfigRouter(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/config", middleware.Auth(), middleware.CheckRole(), config.GetConfig)

	group := apiGroup.Group("/config", middleware.Auth(), middleware.CheckRole())

	group.GET("/limit", middleware.Auth(), middleware.CheckRole(), config.GetConfigLimit)
	group.PATCH("/limit", middleware.Auth(), middleware.CheckRole(), config.UpdateConfigLimit)
}

func TxRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/tx", middleware.Auth(), middleware.CheckRole())
	group.GET("/list", tx.GetTxsList)
}

func CibRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/cib", middleware.Auth(), middleware.CheckRole())
	group.GET("/list", cib.GetCibsList)
	group.POST("/upload", cib.UploadCib)
}

func RiskControlRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/risk_control", middleware.Auth(), middleware.CheckRole())
	group.GET("/roles", riskControl.GetRiskControlRoles)
	group.GET("/:customer_id/role", riskControl.GetCustomerRiskControlRole)
	group.PATCH("/:customer_id/role", riskControl.UpdateCustomerRiskControlRole)
	group.PATCH("/:customer_id/role/reset", riskControl.ResetCustomerRiskControlRole)
	group.PATCH("/role_reset/batch", riskControl.ResetCustomerRiskControlRoleBatch)
	group.PATCH("/:customer_id/limit", riskControl.UpdateCustomerRiskControlLimit)
}

func ViewRouter(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/view", middleware.Auth(), middleware.CheckRole())
	group.GET("/tx_volumns", view.GetTxVolumns)
}

func InternalRouter(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/btm/logs", middleware.Auth(), middleware.CheckRole(), debug.GetBTMChangeLogs)

	group := apiGroup.Group("/btm", middleware.ServerKeyAuth())
	group.GET("/id_number", customer.GetCustomerIdNumber)
	group.POST("/cib", debug.DownlaodCIB)
	group.POST("/add_sumsub_tag", customer.AddSumsubTag)

	toolGroup := apiGroup.Group("/tool", middleware.ServerKeyAuth())
	toolGroup.GET("/sync_sumsub", tool.SyncAllCustomerSumsub)
	toolGroup.GET("/complete_address_binding_log", tool.CompleteAddressBindingLog)
	toolGroup.GET("/sync_daily_tx_volumn", tool.SyncDailyTxVolumn)

}
