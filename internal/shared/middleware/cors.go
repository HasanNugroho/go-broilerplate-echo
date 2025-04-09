package middleware

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetCORS(config *config.Config) echo.MiddlewareFunc {
	allowOrigins := config.Server.AllowedOrigins
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	corsConfig := middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Origin", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	return middleware.CORSWithConfig(corsConfig)
}
