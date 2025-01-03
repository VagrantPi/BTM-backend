package cronjob

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"BTM-backend/third_party/cib"
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func DownlaodCIBAndUpsert() (err error) {
	log := logger.Zap().WithClassFunction("cronjob", "DownlaodCIBAndUpsert")
	defer func() {
		_ = log.Sync()
	}()

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		return
	}

	ctx := context.WithValue(context.Background(), "log", log)

	log.Info("DownlaodCIB")
	config := configs.NewConfigs()

	token, err := cib.GetToken()
	if err != nil {
		log.Error("cib.GetToken()", zap.Any("err", err))
		return
	}

	fmt.Println("token", token)
	zipFile := "cib.zip"
	err = cib.GetWarningZip(token, zipFile)
	if err != nil {
		log.Error("cib.GetWarningZip()", zap.Any("err", err))
		return
	}
	defer os.Remove(zipFile)

	csvFile := "cib.csv"
	err = tools.UnzipFile(zipFile, csvFile, config.Cib.ZipPwd)
	if err != nil {
		log.Error("tools.UnzipFile()", zap.Any("err", err))
		return
	}

	_file, err := os.Open(csvFile)
	if err != nil {
		log.Error("os.Open()", zap.Any("err", err))
		return
	}
	defer _file.Close()
	defer os.Remove(csvFile)

	cibs, err := cib.ConvertCsvFileToBTMCIB(_file)
	if err != nil {
		log.Error("cib.ConvertCsvFileToBTMCIB()", zap.Any("err", err))
		return
	}

	for _, cib := range cibs {
		err = repo.UpsertBTMCIB(repo.GetDb(ctx), cib)
		if err != nil {
			log.Error("repo.UpsertBTMCIB()", zap.Any("err", err))
			continue
		}
	}

	return nil
}
