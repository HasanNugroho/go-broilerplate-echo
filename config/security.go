package config

import (
	"os"

	utils "github.com/HasanNugroho/starter-golang/pkg/utlis"
	"github.com/ulule/limiter/v3"
)

type SecurityConfig struct {
	CheckOrigin     bool
	RateLimit       string
	TrustedPlatform string
	LimiterInstance *limiter.Limiter
}

func LoadSecurityConfig() (securityConfig SecurityConfig) {
	securityConfig.CheckOrigin = utils.ToBool(os.Getenv("ACTIVATE_ORIGIN_VALIDATION"), false)
	securityConfig.RateLimit = utils.ToString(os.Getenv("RATE_LIMIT"), "100-M")
	securityConfig.TrustedPlatform = utils.ToString(os.Getenv("TRUSTED_PLATFORM"), "X-Real-Ip")
	return
}
