package middleware

import (
	"BTM-backend/configs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServerKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		token := c.GetHeader("token")

		// Check if token exists
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Missing server key token in header",
			})
			return
		}

		// Get the server key from environment variable
		serverKey := configs.C.ServerKey

		// Check if SERVER_KEY is set
		if serverKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "SERVER_KEY environment variable is not set",
			})
			return
		}

		// Validate token against server key
		if token != serverKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Invalid server key token",
			})
			return
		}

		// Token is valid, proceed to the next handler
		c.Next()
	}
}
