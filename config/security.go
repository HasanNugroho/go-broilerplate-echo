package config

import (
	"os"

	utils "github.com/HasanNugroho/starter-golang/pkg/utlis"
	"github.com/ulule/limiter/v3"
)

type SecurityConfig struct {
	CheckOrigin       bool
	RateLimit         string
	TrustedPlatform   string
	LimiterInstance   *limiter.Limiter
	ExpectedHost      string
	XFrameOptions     string
	ContentSecurity   string
	XXSSProtection    string
	StrictTransport   string
	ReferrerPolicy    string
	XContentTypeOpts  string
	PermissionsPolicy string
}

func LoadSecurityConfig() (securityConfig SecurityConfig) {
	// Baca setting dari .env
	securityConfig.CheckOrigin = utils.ToBool(os.Getenv("ACTIVATE_ORIGIN_VALIDATION"), false)
	securityConfig.RateLimit = utils.ToString(os.Getenv("RATE_LIMIT"), "100-M")
	securityConfig.TrustedPlatform = utils.ToString(os.Getenv("TRUSTED_PLATFORM"), "X-Real-Ip")
	securityConfig.ExpectedHost = utils.ToString(os.Getenv("EXPECTED_HOST"), "*")
	securityConfig.XFrameOptions = utils.ToString(os.Getenv("X_FRAME_OPTIONS"), "DENY")
	securityConfig.ContentSecurity = utils.ToString(os.Getenv("CONTENT_SECURITY_POLICY"), "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	securityConfig.XXSSProtection = utils.ToString(os.Getenv("X_XSS_PROTECTION"), "1; mode=block")
	securityConfig.StrictTransport = utils.ToString(os.Getenv("STRICT_TRANSPORT_SECURITY"), "max-age=31536000; includeSubDomains; preload")
	securityConfig.ReferrerPolicy = utils.ToString(os.Getenv("REFERRER_POLICY"), "strict-origin")
	securityConfig.XContentTypeOpts = utils.ToString(os.Getenv("X_CONTENT_TYPE_OPTIONS"), "nosniff")
	securityConfig.PermissionsPolicy = utils.ToString(os.Getenv("PERMISSIONS_POLICY"), "")

	return
}
