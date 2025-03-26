package middleware

import (
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		uri := c.Request.URL.Path

		// 轉換為 F2ERouterId
		f2eRouterIds := domain.URIToF2ERouterId(method, uri)
		rolesRaw, _ := c.Get("roles")
		roles, ok := rolesRaw.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "roles not found",
			})
			return
		}

		if len(f2eRouterIds) > 0 && !tools.SliceInSlice(roles, f2eRouterIds) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "forbidden",
			})
			return
		}

		// 繼續處理
		c.Next()
	}
}
