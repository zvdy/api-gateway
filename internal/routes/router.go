package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zvdy/api-gateway/config"
	"github.com/zvdy/api-gateway/internal/auth"
	"github.com/zvdy/api-gateway/internal/proxy"
	"github.com/zvdy/api-gateway/internal/ratelimit"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Rate limiter setup
	err = ratelimit.InitializeRedis(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		panic("Failed to initialize Redis for rate limiting: " + err.Error())
	}
	router.Use(ratelimit.RateLimitMiddleware(100, time.Minute))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Authentication middleware
	router.Use(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token := authHeader[len("Bearer "):]
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	})

	// Reverse proxy setup
	router.Any("/api/*path", proxy.ReverseProxy("http://backend:5000"))

	return router
}
