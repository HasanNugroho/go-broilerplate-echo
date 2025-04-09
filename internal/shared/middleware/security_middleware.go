package middleware

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SecurityMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		Skipper:               middleware.DefaultSkipper,
		XSSProtection:         cfg.Security.XXSSProtection,
		ContentTypeNosniff:    cfg.Security.XContentTypeOpts,
		XFrameOptions:         cfg.Security.XFrameOptions,
		HSTSMaxAge:            63072000, // bisa disesuaikan atau ambil dari config
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    true,
		ContentSecurityPolicy: cfg.Security.ContentSecurity,
		CSPReportOnly:         false,
		ReferrerPolicy:        cfg.Security.ReferrerPolicy,
	})
}
