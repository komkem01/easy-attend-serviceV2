package middlewares

import (
	"easy-attend-service/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Get user ID from context (if authenticated)
		userID := ""
		if uid, exists := c.Get("user_id"); exists {
			userID = uid.(string)
		}

		// Log request
		logger.LogAPIRequest(c.Request.Method, c.Request.URL.Path, userID)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Log response
		logger.Log.WithFields(logrus.Fields{
			"type":        "api_response",
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"latency":     latency.String(),
			"user_id":     userID,
			"ip":          c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		}).Info("API response sent")
	}
}
