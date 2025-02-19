package config

import "github.com/ulule/limiter/v3"

type SecurityConfig struct {
	CheckOrigin     bool
	RateLimit       string
	TrustedPlatform string
	LimiterInstance *limiter.Limiter
}
