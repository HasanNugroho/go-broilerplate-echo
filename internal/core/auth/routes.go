package auth

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, config *config.Config, handler *AuthHandler) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", handler.Login)

		authRoutes.Use(middleware.AuthMiddleware(config))
		authRoutes.POST("/logout", handler.Logout)
		authRoutes.POST("/access-token", handler.GenerateAccessToken)
	}
}
