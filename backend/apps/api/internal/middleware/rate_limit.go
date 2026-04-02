package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	Client *redis.Client
	Limit  int64
	Window time.Duration
}

func NewRateLimiter(client *redis.Client, limit int64, window time.Duration) *RateLimiter {
	return &RateLimiter{
		Client: client,
		Limit:  limit,
		Window: window,
	}
}

func (r *RateLimiter) Middleware(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if r == nil || r.Client == nil {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
		defer cancel()

		key := fmt.Sprintf("rl:%s:%s", scope, c.ClientIP())
		current, err := r.Client.Incr(ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "rate limiter unavailable"})
			return
		}

		if current == 1 {
			if err := r.Client.Expire(ctx, key, r.Window).Err(); err != nil {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "rate limiter unavailable"})
				return
			}
		}

		if current > r.Limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}

		c.Next()
	}
}
