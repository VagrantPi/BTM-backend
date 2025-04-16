package middleware

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"encoding/json"
	"fmt"
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

		userInfo, err := tools.ParseToken(token, configs.C.JWT.Secret)
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
		role, ok := domain.GetTTLMap[domain.RoleWithTTL](&domain.TTLRoleMap, fmt.Sprintf("%d", userInfo.Role))
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

			expire := time.Now().Add(30 * time.Minute).UnixNano()
			// cache
			roleWithTTL := domain.TTLMap[domain.RoleWithTTL]{
				Cache: domain.RoleWithTTL{
					CacheRole:  fetchRole,
					Expiration: expire,
				},
				Expire: expire,
			}

			// 使用正確的結構調用 SetTTLMap
			domain.SetTTLMap[domain.RoleWithTTL](&domain.TTLRoleMap, fmt.Sprintf("%d", userInfo.Role), roleWithTTL.Cache, roleWithTTL.Expire)
			role = &roleWithTTL.Cache
		}

		var roles []domain.RoleItem
		_ = json.Unmarshal([]byte(role.CacheRole.RoleRaw), &roles)
		c.Set("roles", domain.PageIds(roles))
		c.Set("role", role.CacheRole)

		if configs.C.Env != "local" {
			// 只允許一個裝置登入
			isLastLoginToken, err := repo.IsLastLoginToken(repo.GetDb(c), uint(userInfo.Id), token)
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
		}

		c.Set("userInfo", userInfo)

		c.Next()
	}
}
