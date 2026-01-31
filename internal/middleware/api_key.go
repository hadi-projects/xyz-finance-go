package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientKey := c.GetHeader("X-API-KEY")
		if clientKey == "" || clientKey != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid API Key"})
			c.Abort()
			return
		}
		c.Next()
	}
}
