package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDContextKey = "request_id"

func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if logger == nil {
			c.Next()
			return
		}

		start := time.Now()
		requestID := uuid.NewString()
		c.Set(requestIDContextKey, requestID)

		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()

		logger.Info(
			"http_request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

func RequestID(c *gin.Context) (string, bool) {
	raw, ok := c.Get(requestIDContextKey)
	if !ok {
		return "", false
	}
	id, ok := raw.(string)
	return id, ok
}
