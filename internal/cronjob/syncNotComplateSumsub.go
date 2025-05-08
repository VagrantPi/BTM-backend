package cronjob

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/logger"
	"BTM-backend/third_party/sumsub"
	"context"
	"time"

	"go.uber.org/zap"
)

func SyncNotComplateSumsub(force bool) {
	log := logger.Zap().WithClassFunction("cronjob", "SyncNotComplateSumsub")
	defer func() {
		_ = log.Sync()
	}()

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		return
	}

	ctx := context.Background()
	ids, err := repo.GetUnCompletedSumsubCustomerIds(repo.GetDb(ctx), force)
	if err != nil {
		log.Error("repo.GetUnCompletedSumsubCustomerIds()", zap.Any("err", err))
		return
	}
	log.Info("開始每日 sync Sumsub 資料", zap.Int("資料筆數", len(ids)))

	for _, id := range ids {
		log.Info("-----------------------------------------")
		log.Info("sync Sumsub 資料", zap.String("customer_id", id))

		sumsub.FetchDataAdapter(ctx, log, repo, id)
		time.Sleep(500 * time.Microsecond)
	}

}
