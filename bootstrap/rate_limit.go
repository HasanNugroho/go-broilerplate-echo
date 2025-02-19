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
func InitRateLimiter(formattedRateLimit, trustedPlatform string) (*limiter.Limiter, error) {
	cfg := config.GetConfig()

	if formattedRateLimit == "" {
		return nil, nil
	}

	rate, err := limiter.NewRateFromFormatted(formattedRateLimit)
	if err != nil {
		return nil, err
	}

	var limiterInstance *limiter.Limiter
	// custom IPv6 mask
	ipv6Mask := net.CIDRMask(64, 128)

	options := []limiter.Option{
		limiter.WithIPv6Mask(ipv6Mask),
	}

	if cfg.Database.REDIS.Activate {
		// Create a store with the redis client.
		store, err := sredis.NewStoreWithOptions(cfg.Database.REDIS.Client, limiter.StoreOptions{
			Prefix:   "limiter",
			MaxRetry: 3,
		})
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		limiterInstance = limiter.New(
			store,
			rate,
			options...,
		)
	} else {
		// use an in-memory store with a goroutine which clears expired keys
		store := memory.NewStore()

		if trustedPlatform != "" {
			options = append(options, limiter.WithClientIPHeader(trustedPlatform))
		}

		// create the limiter instance
		limiterInstance = limiter.New(
			store,
			rate,
			options...,
		)
	}

	cfg.Security.LimiterInstance = limiterInstance

	return limiterInstance, nil
}
