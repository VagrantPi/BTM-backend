package cronjob

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/logger"
	"context"

	"go.uber.org/zap"
)

func ShapshotYesterdayVolumn() (err error) {
	log := logger.Zap().WithClassFunction("cronjob", "ShapshotVolumn")
	defer func() {
		_ = log.Sync()
	}()

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		return
	}

	ctx := context.WithValue(context.Background(), "log", log)

	log.Info("ShapshotVolumn")

	return repo.SnapshotYesterday(repo.GetDb(ctx))
}
