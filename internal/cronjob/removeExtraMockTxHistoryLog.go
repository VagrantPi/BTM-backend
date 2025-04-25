package cronjob

import (
	"BTM-backend/internal/di"
	"context"
)

func RemoveExtraMockTxHistoryLog() (err error) {
	repo, err := di.NewRepo()
	if err != nil {
		return
	}

	return repo.RemoveExtraMockTxHistoryLog(repo.GetDb(context.Background()))
}
