package bootstrap

import (
	"log"
	"net"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// InitRateLimiter - initialize the rate limiter instance
func InitRateLimiter(cfg *config.Configuration, formattedRateLimit string, trustedPlatform string) (limiterInstance *limiter.Limiter, err error) {
	if formattedRateLimit == "" {
		return nil, nil
	}

	rate, err := limiter.NewRateFromFormatted(formattedRateLimit)
	if err != nil {
		return nil, err
	}

	ipv6Mask := net.CIDRMask(64, 128)
	options := []limiter.Option{limiter.WithIPv6Mask(ipv6Mask)}

	if cfg.Database.Redis.Enabled {
		// Create a store with the redis client.
		if cfg.Database.Redis.Client == nil {
			log.Fatal("Redis client is not initialized")
		}

		store, err := sredis.NewStoreWithOptions(cfg.Database.Redis.Client, limiter.StoreOptions{
			Prefix:   "limiter",
			MaxRetry: 3,
		})

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		cfg.Security.LimiterInstance = limiter.New(store, rate, options...)
		return cfg.Security.LimiterInstance, nil
	}

	// default use memory store
	store := memory.NewStore()
	if trustedPlatform != "" {
		options = append(options, limiter.WithClientIPHeader(trustedPlatform))
	}

	cfg.Security.LimiterInstance = limiter.New(store, rate, options...)
	return cfg.Security.LimiterInstance, nil
}
