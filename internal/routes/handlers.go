package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zvdy/api-gateway/internal/ratelimit"
)

// RateLimitMiddleware applies rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, err := ratelimit.AllowRequest(c, c.ClientIP(), 100, time.Minute)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "rate limiting error"})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}

// ProxyHandler handles the reverse proxy
func ProxyHandler(c *gin.Context) {
	// Implement reverse proxy logic here
	c.JSON(http.StatusOK, gin.H{"message": "proxying request"})
}
