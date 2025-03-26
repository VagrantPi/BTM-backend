package middleware

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"encoding/json"
	"net/http"
	"time"

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

		// 取得剛用戶權限路由表
		role, ok := domain.GetTTLRoleMap(uint(userInfo.Role))
		if !ok {
			fetchRole, err := repo.GetRawRoleById(repo.GetDb(c), userInfo.Role)
			if err != nil {
				log.Error("GetRawRoleById", zap.Any("err", err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "GetRawRoleById error",
				})
				return
			}
			// TTL 不宜太長，因為前端可以動態更動
			domain.SetTTLRoleMap(uint(userInfo.Role), fetchRole, time.Now().Add(1*time.Minute).UnixNano())
			role = &fetchRole
		}

		var roles []domain.RoleItem
		_ = json.Unmarshal([]byte(role.RoleRaw), &roles)
		c.Set("roles", domain.PageIds(roles))

		// 只允許一個裝置登入
		isLastLoginToken, err := repo.IsLastLoginToken(repo.GetDb(c), userInfo.Id, token)
		if err != nil {
			log.Error("IsLastLoginToken error", zap.Any("err", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "IsLastLoginToken error",
			})
			return
		}

		if !isLastLoginToken {
			log.Error("login only on device", zap.Any("userId", userInfo.Id))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "只能登入單一裝置，請先登出其他裝置在進行登入",
			})
			return
		}

		c.Set("userInfo", userInfo)

		c.Next()
	}
}
