package main

import (
	"BTM-backend/configs"
	"BTM-backend/internal/controller"
	"BTM-backend/internal/cronjob"
	"BTM-backend/internal/middleware"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/logger"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	log := logger.Zap()
	defer log.Sync()

	// 設置 Gin 的日誌格式
	gin.DefaultWriter = io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename: "logs/app.log",
		MaxAge:   1825,
		Compress: true,
	})

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
	log.Info("Application starting...")

	if err := s.ListenAndServe(); err != nil {
		log.Error("s.ListenAndServe()", zap.Any("err", err))
		panic(err)
	}
}

// 自定義的 Gin 日誌格式
func customGinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logger := logger.Zap()
		logger.Info("Gin request",
			zap.Any("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("duration", duration),
		)
	}
}

func setupGin() http.Handler {

	mode := configs.C.Gin.Mode
	gin.SetMode(mode)

	r := gin.New()

	r.Use(customGinLogger())
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
	controller.SystemRouter(apiGroup)

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

	// 每天處理未完成的 Sumsub 資料 - 每日 1:30
	cron.AddFunc("30 1 * * *", func() {
		cronjob.SyncNotComplateSumsub(false)
	})

	// 每天凌晨 12:30 快照前天交易量 by 個台機器
	cron.AddFunc("30 0 * * *", func() {
		cronjob.ShapshotYesterdayVolumn()
	})

	// 每個小時，移除多餘的新增限額塞入假資料 log
	cron.AddFunc("0 * * * *", func() {
		cronjob.RemoveExtraMockTxHistoryLog()
	})

	cron.Start()
}
