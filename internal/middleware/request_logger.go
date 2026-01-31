package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hadi-projects/xyz-finance-go/pkg/logger"
	"github.com/rs/zerolog"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Generate Request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Process Request
		c.Next()

		// Statistics
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		userID, userExists := c.Get("user_id")

		// Determine Logger
		var logEvent *zerolog.Event
		// Assuming Auth routes are prefixed with /api/auth
		if strings.HasPrefix(path, "/api/auth") {
			if statusCode >= 500 {
				logEvent = logger.AuthLogger.Error()
			} else {
				logEvent = logger.AuthLogger.Info()
			}
		} else {
			if statusCode >= 500 {
				logEvent = logger.SystemLogger.Error()
			} else {
				logEvent = logger.SystemLogger.Info()
			}
		}

		// Build Log Entry
		logEvent.
			Str("request_id", requestID).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Str("ip", clientIP).
			Str("user_agent", userAgent).
			Dur("latency", latency)

		if userExists {
			logEvent.Uint("user_id", userID.(uint))
		}

		// Log Errors if any
		if len(c.Errors) > 0 {
			logEvent.Str("error", c.Errors.String())
		}

		logEvent.Msg("handled request")
	}
}
