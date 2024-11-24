package ratelimit

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitializeRedis(addr, password string, db int) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	return err
}

func AllowRequest(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	// Use Redis to implement rate limiting
	pipe := RedisClient.TxPipeline()
	defer pipe.Close()

	// Increment the counter for the given key
	count := pipe.Incr(ctx, key)
	// Set the expiration for the key if it doesn't exist
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// Check if the request count exceeds the limit
	if count.Val() > int64(limit) {
		return false, nil
	}

	return true, nil
}
