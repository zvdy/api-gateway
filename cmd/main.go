package main

import (
	"log"

	"github.com/zvdy/api-gateway/config"
	"github.com/zvdy/api-gateway/internal/ratelimit"
	"github.com/zvdy/api-gateway/internal/routes"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Redis for rate limiting
	ratelimit.InitializeRedis(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)

	// Set up router
	router := routes.SetupRouter()

	log.Printf("API Gateway is running at %s\n", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
