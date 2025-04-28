package tools

import (
	"BTM-backend/internal/controller"
	"BTM-backend/internal/middleware"
	"BTM-backend/pkg/api"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupGin() http.Handler {
	gin.SetMode("debug")

	r := gin.New()

	r.Use(middleware.CORS)
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(middleware.RequestId)
	r.Use(middleware.ErrHandler)

	apiGroup := r.Group("/api")

	apiGroup.GET("/ping", api.Ping)

	// 以下 API 以 /api 為前綴
	controller.UserRouter(apiGroup)
	controller.ThirdPartyRouter(apiGroup)
	controller.CustomerRouter(apiGroup)
	controller.InternalRouter(apiGroup)
	controller.UserConfigRouter(apiGroup)
	controller.TxRouter(apiGroup)
	controller.CibRouter(apiGroup)
	controller.RiskControlRouter(apiGroup)
	controller.ViewRouter(apiGroup)

	return r
}
