package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BlockDestructiveWhenImpersonating() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := strings.TrimSpace(c.GetHeader("X-Impersonation-Session-ID"))
		if sessionID == "" {
			c.Next()
			return
		}

		method := strings.ToUpper(c.Request.Method)
		if method == http.MethodDelete {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "destructive actions are blocked during impersonation",
			})
			return
		}
		c.Next()
	}
}
