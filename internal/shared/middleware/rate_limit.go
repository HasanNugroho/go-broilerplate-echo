package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func RateLimit(config *config.Config) echo.MiddlewareFunc {
	cfg := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(config.Security.RateLimit),
				Burst:     30,
				ExpiresIn: 3 * time.Minute, // TTL untuk token
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil // bisa diganti pakai user ID dari JWT juga
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "Forbidden",
			})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, echo.Map{
				"message": "Too many requests",
			})
		},
	}

	fmt.Println("✅ RateLimiter middleware applied")
	return middleware.RateLimiterWithConfig(cfg)
}

// // RateLimit - rate limit middleware
// func RateLimit(limiterInstance *limiter.Limiter) echo.MiddlewareFunc {
// 	// Kalau limiter-nya nil
// 	if limiterInstance == nil {
// 		fmt.Println("⚠️ Limiter instance is nil, skipping middleware")
// 		return func(next echo.HandlerFunc) echo.HandlerFunc {
// 			return func(c echo.Context) error {
// 				return next(c)
// 			}
// 		}
// 	}

// 	// Middleware aktif
// 	fmt.Println("✅ RateLimit middleware applied")

// 	// Inisialisasi middleware Echo pakai limiter instance
// 	middleware := mecho.NewMiddleware(limiterInstance)

// 	return middleware
// }

// InitRateLimiter - initialize the rate limiter instance
// func InitRateLimiter(config *config.Config, redisClient *redisClient.Client, formattedRateLimit string, trustedPlatform string) (limiterInstance *limiter.Limiter, err error) {
// 	if formattedRateLimit == "" {
// 		return nil, nil
// 	}

// 	rate, err := limiter.NewRateFromFormatted(formattedRateLimit)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ipv6Mask := net.CIDRMask(64, 128)
// 	options := []limiter.Option{limiter.WithIPv6Mask(ipv6Mask)}

// 	if config.Redis.Enabled {
// 		// Create a store with the redis client.
// 		if redisClient == nil {
// 			log.Fatal("Redis client is not initialized")
// 		}

// 		store, err := redis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
// 			Prefix:   "limiter",
// 			MaxRetry: 3,
// 		})

// 		if err != nil {
// 			log.Fatal(err)
// 			return nil, err
// 		}

// 		config.Security.LimiterInstance = limiter.New(store, rate, options...)
// 		return config.Security.LimiterInstance, nil
// 	}

// 	// default use memory store
// 	store := memory.NewStore()
// 	if trustedPlatform != "" {
// 		options = append(options, limiter.WithClientIPHeader(trustedPlatform))
// 	}

// 	config.Security.LimiterInstance = limiter.New(store, rate, options...)
// 	return config.Security.LimiterInstance, nil
// }
