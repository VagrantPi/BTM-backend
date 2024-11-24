package main

import (
	"BTM-backend/configs"
	"BTM-backend/internal/controller"
	"BTM-backend/internal/middleware"
	"BTM-backend/pkg/api"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	h2 := h2c.NewHandler(setupGin(), &http2.Server{})
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", configs.C.Port),
		Handler:        h2,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   configs.C.TimeoutSecond * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("開始監聽 %v\n", configs.C.Port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func setupGin() http.Handler {

	mode := configs.C.Gin.Mode
	gin.SetMode(mode)

	r := gin.New()

	// 以下為一個請求流

	r.Use(middleware.CORS)
	r.Use(middleware.RequestId)
	r.Use(middleware.ErrHandler)

	apiGroup := r.Group("/api")

	apiGroup.GET("/ping", api.Ping)

	// apiGroup.POST("/ping", pkg.Wrap(&di.ApiPing))

	// 以下 API 以 /api 為前綴
	// routeClient(apiGroup)
	controller.UserRouter(apiGroup)

	return r
}
