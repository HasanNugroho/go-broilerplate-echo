package configs

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/HasanNugroho/starter-golang/internal/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
	ConnTTL  int
	Client   *redis.Client
}

// LoadRedisConfig loads Redis configuration
func loadRedisConfig() (redisConfig RedisConfig) {
	redisConfig = RedisConfig{
		Enabled:  utils.ToBool(os.Getenv("ACTIVATE_REDIS"), false),
		Host:     utils.ToString(os.Getenv("REDISHOST"), "localhost"),
		Port:     utils.ToInt(os.Getenv("REDISPORT"), 6379),
		Password: utils.ToString(os.Getenv("REDISPASSWORD"), ""),
		PoolSize: utils.ToInt(os.Getenv("POOLSIZE"), 10),
		ConnTTL:  utils.ToInt(os.Getenv("CONNTTL"), 60),
	}
	return
}

// InitRedis initializes the Redis client
func InitRedis(cfg *RedisConfig) (*redis.Client, error) {
	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port)),
		Password: cfg.Password,
		PoolSize: cfg.PoolSize,
	})

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ConnTTL)*time.Second)
	defer cancel()

	// Test Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("❌ Redis connection failed: %w", err)
	}

	// Assign Redis client to config
	cfg.Client = redisClient
	Logger.Info().Msg("✅ Redis connected successfully!")

	return redisClient, nil
}

// ShutdownRedis closes the Redis connection safely
func ShutdownRedis(redisClient *redis.Client) {
	if redisClient == nil {
		Logger.Info().Msg("⚠️ Redis client is nil, skipping shutdown.")
		return
	}

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Close(); err != nil {
		Logger.Error().Msgf("❌ Error closing Redis connection: %v", err)
	} else {
		Logger.Info().Msg("✅ Redis connection closed successfully!")
	}
}
