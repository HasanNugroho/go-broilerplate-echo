package users

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(router *gin.RouterGroup, handler *UserHandler) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", handler.Create)
		userRoutes.GET("/", handler.FindAll)
		userRoutes.GET("/:id", handler.FindById)
		userRoutes.PUT("/:id", handler.Update)
		userRoutes.DELETE("/:id", handler.Delete)
	}
}
