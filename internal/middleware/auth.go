package middleware

import (
	"BTM-backend/internal/di"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Auth() gin.HandlerFunc {
	log := logger.Zap().WithClassFunction("middleware", "Auth")
	defer func() {
		_ = log.Sync()
	}()

	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			log.Error("token is empty")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "token is empty",
			})
			return
		}

		userInfo, err := tools.ParseToken(token)
		if err != nil {
			log.Error("parseToken error", zap.Any("err", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "parseToken error",
			})
			return
		}

		repo, err := di.NewRepo()
		if err != nil {
			log.Error("di repo error", zap.Any("err", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "di repo error",
			})
			return
		}

		// isLastLoginToken, err := repo.IsLastLoginToken(repo.GetDb(c), userInfo.Id, token)
		// if err != nil {
		// 	log.Error("IsLastLoginToken error", zap.Any("err", err))
		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		// 		"code": http.StatusInternalServerError,
		// 		"msg":  "IsLastLoginToken error",
		// 	})
		// 	return
		// }

		// if !isLastLoginToken {
		// 	log.Error("login only on device", zap.Any("userId", userInfo.Id))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		// 		"code": http.StatusForbidden,
		// 		"msg":  "login only on device",
		// 	})
		// 	return
		// }

		c.Set("userInfo", userInfo)

		c.Next()
	}
}
