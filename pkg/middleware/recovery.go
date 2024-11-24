package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zvdy/api-gateway/internal/auth"
	"github.com/zvdy/api-gateway/internal/proxy"
	"github.com/zvdy/api-gateway/internal/ratelimit"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS setup
	router.Use(corsMiddleware())

	// Rate limiter setup
	err := ratelimit.InitializeRedis("localhost:6379", "", 0)
	if err != nil {
		panic("Failed to initialize Redis for rate limiting: " + err.Error())
	}
	router.Use(RateLimitMiddleware(100, time.Minute))

	// Authentication middleware
	router.Use(AuthMiddleware())

	// Reverse proxy setup
	router.Any("/api/*path", proxy.ReverseProxy("http://localhost:5000"))

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, err := ratelimit.AllowRequest(c, c.ClientIP(), limit, window)
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		_, err := auth.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Next()
	}
}
