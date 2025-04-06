package config

import (
	"log"
	"net"

	redisClient "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	redis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// InitRateLimiter - initialize the rate limiter instance
func InitRateLimiter(config *Config, redisClient *redisClient.Client, formattedRateLimit string, trustedPlatform string) (limiterInstance *limiter.Limiter, err error) {
	if formattedRateLimit == "" {
		return nil, nil
	}

	rate, err := limiter.NewRateFromFormatted(formattedRateLimit)
	if err != nil {
		return nil, err
	}

	ipv6Mask := net.CIDRMask(64, 128)
	options := []limiter.Option{limiter.WithIPv6Mask(ipv6Mask)}

	if config.Redis.Enabled {
		// Create a store with the redis client.
		if redisClient == nil {
			log.Fatal("Redis client is not initialized")
		}

		store, err := redis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
			Prefix:   "limiter",
			MaxRetry: 3,
		})

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		config.Security.LimiterInstance = limiter.New(store, rate, options...)
		return config.Security.LimiterInstance, nil
	}

	// default use memory store
	store := memory.NewStore()
	if trustedPlatform != "" {
		options = append(options, limiter.WithClientIPHeader(trustedPlatform))
	}

	config.Security.LimiterInstance = limiter.New(store, rate, options...)
	return config.Security.LimiterInstance, nil
}
