package auth

import "github.com/gin-gonic/gin"

func RegisterAuthRoutes(router *gin.RouterGroup, handler *AuthHandler) {
	userRoutes := router.Group("/auth")
	{
		userRoutes.POST("/login", handler.Login)
	}
}
