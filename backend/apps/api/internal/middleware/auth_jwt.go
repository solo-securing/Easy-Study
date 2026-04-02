package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const claimsContextKey = "jwt_claims"

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.TrimSpace(secret) == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "jwt secret not configured"})
			return
		}

		authorization := c.GetHeader("Authorization")
		tokenString, ok := strings.CutPrefix(authorization, "Bearer ")
		if !ok || strings.TrimSpace(tokenString) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		c.Set(claimsContextKey, claims)
		c.Next()
	}
}

func JWTClaims(c *gin.Context) (jwt.MapClaims, bool) {
	raw, ok := c.Get(claimsContextKey)
	if !ok {
		return nil, false
	}

	claims, ok := raw.(jwt.MapClaims)
	return claims, ok
}
