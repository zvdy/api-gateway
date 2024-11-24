package cache

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedisClient(redisAddr string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}

func RateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		pipe := redisClient.Pipeline()
		pipe.Incr(c, key)
		pipe.Expire(c, key, window)
		result, err := pipe.Exec(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		count := result[0].(*redis.IntCmd).Val()
		if count > int64(maxRequests) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
