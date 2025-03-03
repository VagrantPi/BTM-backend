package main

import (
	"BTM-backend/configs"
	"BTM-backend/internal/controller"
	"BTM-backend/internal/cronjob"
	"BTM-backend/internal/middleware"
	"BTM-backend/pkg/api"
	"fmt"
	"log"
	"net/http"
	"time"
	_ "time/tzdata"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	// 背景工作
	go BackgroundWorker()

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

	r.Use(middleware.CORS)
	r.Use(middleware.RequestId)
	r.Use(middleware.ErrHandler)

	apiGroup := r.Group("/api")

	apiGroup.GET("/ping", api.Ping)

	// 後台
	r.Static("/admin", "./dist")

	// 以下 API 以 /api 為前綴
	controller.UserRouter(apiGroup)
	controller.CustomerRouter(apiGroup)
	controller.CustomerInternalRouter(apiGroup)
	controller.UserConfigRouter(apiGroup)
	controller.TxRouter(apiGroup)
	controller.DebugRouter(apiGroup)
	controller.CibRouter(apiGroup)

	return r
}

func BackgroundWorker() {
	CronJob()
}

func CronJob() {
	fmt.Println("開始定時任務")
	// 使用的時區
	nyc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
	// 宣告可以使用秒
	cron := cron.New(cron.WithSeconds(), cron.WithLocation(nyc))

	// 每天處理告誡名單 - 每日 1:00
	cron.AddFunc("0 1 * * *", func() {
		cronjob.DownlaodCIBAndUpsert()
	})

	cron.Start()
}
