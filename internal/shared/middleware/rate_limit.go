package middleware

import (
	"fmt"
	"log"
	"net"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"

	redisClient "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	redis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RateLimit - rate limit middleware
func RateLimit(limiterInstance *limiter.Limiter) gin.HandlerFunc {
	// limiter instance is nil
	if limiterInstance == nil {
		fmt.Println("⚠️ Limiter instance is nil, skipping middleware")
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// middleware aktif
	fmt.Println("✅ RateLimit middleware applied")
	// give the limiter instance to the middleware initializer
	return mgin.NewMiddleware(limiterInstance)
}

// InitRateLimiter - initialize the rate limiter instance
func InitRateLimiter(config *config.Config, redisClient *redisClient.Client, formattedRateLimit string, trustedPlatform string) (limiterInstance *limiter.Limiter, err error) {
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
