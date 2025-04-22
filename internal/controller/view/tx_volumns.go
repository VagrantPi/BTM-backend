package view

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetTxVolumnsRep struct {
	Today           []domain.DeviceData `json:"today"`
	TodayTotal      int64               `json:"today_total"`
	SevenDays       []domain.DeviceData `json:"seven_days"`
	SevenDaysTotal  int64               `json:"seven_days_total"`
	ThirtyDays      []domain.DeviceData `json:"thirty_days"`
	ThirtyDaysTotal int64               `json:"thirty_days_total"`
}

func GetTxVolumns(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetTxVolumns")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	rep := GetTxVolumnsRep{}
	mux := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(3)

	// 今日
	go func() {
		defer wg.Done()

		today := time.Now().Format("2006-01-02")
		txs, total, err := repo.FetchByStatDateAndGroupByDeviceId(repo.GetDb(c), today, today)
		if err != nil {
			log.Error("repo.GetCustomers()", zap.Any("err", err))
		}

		mux.Lock()
		rep.Today = txs
		rep.TodayTotal = total
		mux.Unlock()
	}()

	// 七天
	go func() {
		defer wg.Done()

		today := time.Now().Format("2006-01-02")
		fmt.Println("today", today)
		start := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		fmt.Println("start", start)
		txs, total, err := repo.FetchByStatDateAndGroupByDeviceId(repo.GetDb(c), start, today)
		if err != nil {
			log.Error("repo.GetCustomers()", zap.Any("err", err))
		}

		mux.Lock()
		rep.SevenDays = txs
		rep.SevenDaysTotal = total
		mux.Unlock()
	}()

	// 30天
	go func() {
		defer wg.Done()

		today := time.Now().Format("2006-01-02")
		start := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		txs, total, err := repo.FetchByStatDateAndGroupByDeviceId(repo.GetDb(c), start, today)
		if err != nil {
			log.Error("repo.GetCustomers()", zap.Any("err", err))
		}

		mux.Lock()
		rep.ThirtyDays = txs
		rep.ThirtyDaysTotal = total
		mux.Unlock()
	}()

	wg.Wait()

	api.OKResponse(c, rep)
}
