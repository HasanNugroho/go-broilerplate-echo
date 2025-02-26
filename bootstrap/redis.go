package bootstrap

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/redis/go-redis/v9"
)

// InitRedis initializes the Redis client
func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Env.Host, strconv.Itoa(cfg.Env.Port)),
		Password: cfg.Env.Password,
		PoolSize: cfg.Conn.PoolSize,
	})

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Conn.ConnTTL)*time.Second)
	defer cancel()

	// Test Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("❌ Redis connection failed: %w", err)
	}

	// Assign Redis client to config
	cfg.Client = redisClient
	log.Println("✅ Redis connected successfully!")

	return redisClient, nil
}

// ShutdownRedis closes the Redis connection safely
func ShutdownRedis(redisClient *redis.Client) {
	if redisClient == nil {
		log.Println("⚠️ Redis client is nil, skipping shutdown.")
		return
	}

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Close(); err != nil {
		log.Printf("❌ Error closing Redis connection: %v", err)
	} else {
		log.Println("✅ Redis connection closed successfully!")
	}
}
