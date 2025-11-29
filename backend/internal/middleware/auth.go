package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1002,
				"message": "missing authorization header",
			})
			c.Abort()
			return
		}

		// 支持 Bearer token 格式
		token := strings.TrimPrefix(auth, "Bearer ")
		if token != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1002,
				"message": "invalid api key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
