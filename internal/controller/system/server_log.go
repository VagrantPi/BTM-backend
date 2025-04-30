package system

import (
	"BTM-backend/pkg/api"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type GetServerLogsReq struct {
	Limit   int       `form:"limit" binding:"required"`
	Page    int       `form:"page" binding:"required"`
	StartAt time.Time `form:"start_at"`
	EndAt   time.Time `form:"end_at"`
}

type GetServerLogsData struct {
	Total int      `json:"total"`
	Items []string `json:"items"`
}

func GetServerLogs(c *gin.Context) {
	req := GetServerLogsReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logs, err := readLogFiles("logs", req.StartAt, req.EndAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	start := (req.Page - 1) * req.Limit
	end := start + req.Limit
	if start > len(logs) {
		start = len(logs)
	}
	if end > len(logs) {
		end = len(logs)
	}

	api.OKResponse(c, GetServerLogsData{
		Total: len(logs),
		Items: logs[start:end],
	})
}

func readLogFiles(dirPath string, startAt, endAt time.Time) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var logs []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if name == ".gitkeep" {
			continue
		}
		if startAt.IsZero() && endAt.IsZero() {
			logs = append(logs, name)
			continue
		}
		t, err := time.Parse("app-2006-01-02T15-04-05.000.log.gz", name)
		if err != nil {
			continue
		}
		if !startAt.IsZero() && t.Before(startAt) {
			continue
		}
		if !endAt.IsZero() && t.After(endAt) {
			continue
		}
		if name == "app.log" {
			logs = append(logs, name)
			continue
		}
		logs = append(logs, name)
	}
	sort.Slice(logs, func(i, j int) bool {
		return logs[i] > logs[j]
	})
	return logs, nil
}
