package system

import (
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

// DownloadLog handles GET /system/download_log?name=xxxx
func DownloadLog(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		api.ErrResponse(c, "c.Query(\"name\")", errors.BadRequest(error_code.ErrInvalidRequest, "c.Query(\"name\")"))
		return
	}

	filePath := filepath.Join("logs", name)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("file not found: %s", filePath),
		})
		return
	}

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	c.Writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	http.ServeFile(c.Writer, c.Request, filePath)
}
