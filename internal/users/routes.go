package users

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(router *gin.RouterGroup, handler *UserHandler) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", handler.Create)
		userRoutes.GET("/", handler.FindAll)
	}
}
