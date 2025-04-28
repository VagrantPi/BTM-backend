package cronjob

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"context"
)

func RemoveExtraMockTxHistoryLog() (err error) {
	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		return
	}

	return repo.RemoveExtraMockTxHistoryLog(repo.GetDb(context.Background()))
}
