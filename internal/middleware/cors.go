package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
)

// CORS 資安, 跨來源資源共用
var CORS = cors.New(
	cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	},
)
