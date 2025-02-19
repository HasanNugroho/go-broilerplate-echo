package middleware

import (
	"github.com/HasanNugroho/starter-golang/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetCORS() gin.HandlerFunc {
	cfg := config.GetConfig()

	allowOrigins := cfg.Server.AllowedOrigins
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	corsConfig := cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Origin", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	return cors.New(corsConfig)
}
