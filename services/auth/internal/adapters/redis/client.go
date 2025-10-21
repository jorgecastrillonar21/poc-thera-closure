package redis

import (
	"theraclosure/auth-service/internal/adapters/config"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return client, nil
}
