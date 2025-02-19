package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/redis/go-redis/v9"
)

// InitRedis initializes the Redis client
func InitRedis() (*redis.Client, error) {
	cfg := config.GetConfig().Database.REDIS

	if cfg.Env.Host == "" || cfg.Env.Port == "" {
		return nil, fmt.Errorf("❌ Redis configuration is missing host or port")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Env.Host, cfg.Env.Port),
		PoolSize: cfg.Conn.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Conn.ConnTTL)*time.Second)
	defer cancel()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("❌ Redis connection failed: %w", err)
	}

	log.Println("✅ Redis connected successfully!")
	return redisClient, nil
}

// ShutdownRedis closes the Redis connection
func ShutdownRedis(redisClient *redis.Client) {
	if redisClient != nil {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := redisClient.Close(); err != nil {
			log.Printf("❌ Error closing Redis connection: %v", err)
		} else {
			log.Println("✅ Redis connection closed successfully!")
		}
	}
}
