package cib

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/logger"
	"BTM-backend/third_party/cib"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UploadCib(c *gin.Context) {
	file, _ := c.FormFile("file")

	log := logger.Zap().WithClassFunction("api", "UploadCib")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		c.JSON(http.StatusOK, gin.H{
			"message": "server error",
		})
		return
	}

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, "./"+file.Filename)

	_file, err := os.Open("./" + file.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "file open error",
		})
		return
	}
	defer _file.Close()
	defer os.Remove("./" + file.Filename)

	cibs, err := cib.ConvertCsvFileToBTMCIB(_file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "file parse error",
		})
		return
	}
	for _, cib := range cibs {
		if cib.DataType == "D" {
			err = repo.DeleteBTMCIB(repo.GetDb(c), cib.Pid)
			if err != nil {
				log.Error("repo.DeleteBTMCIB()", zap.Any("err", err))
				continue
			}
		} else {
			err = repo.UpsertBTMCIB(repo.GetDb(c), cib)
			if err != nil {
				log.Error("repo.UpsertBTMCIB()", zap.Any("err", err))
				continue
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
