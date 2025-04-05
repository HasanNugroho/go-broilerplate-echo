package users

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, config *config.Config, handler *UserHandler) {
	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(config))
	{
		userRoutes.POST("/", handler.Create)
		userRoutes.GET("/", handler.FindAll)
		userRoutes.GET("/:id", handler.FindById)
		userRoutes.PUT("/:id", handler.Update)
		userRoutes.DELETE("/:id", handler.Delete)
	}

}
