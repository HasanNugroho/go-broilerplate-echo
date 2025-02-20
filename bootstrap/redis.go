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
	cfg := config.GetConfig() // Get global configuration
	if cfg == nil {
		return nil, fmt.Errorf("❌ Failed to load config: config is nil")
	}

	redisCfg := cfg.Database.REDIS

	// Validate Redis configuration
	if redisCfg.Env.Host == "" || redisCfg.Env.Port == "" {
		return nil, fmt.Errorf("❌ Redis configuration is missing host or port")
	}

	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisCfg.Env.Host, redisCfg.Env.Port),
		Password: redisCfg.Env.Password,
		PoolSize: redisCfg.Conn.PoolSize,
	})

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(redisCfg.Conn.ConnTTL)*time.Second)
	defer cancel()

	// Test Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("❌ Redis connection failed: %w", err)
	}

	// Assign Redis client to config
	redisCfg.Client = redisClient
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
