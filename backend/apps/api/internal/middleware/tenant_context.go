package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const tenantContextKey = "tenant_id"

func TenantContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "missing tenant context",
			})
			return
		}

		c.Set(tenantContextKey, tenantID)
		c.Next()
	}
}

func TenantID(c *gin.Context) (string, bool) {
	raw, ok := c.Get(tenantContextKey)
	if !ok {
		return "", false
	}
	value, ok := raw.(string)
	return value, ok
}
