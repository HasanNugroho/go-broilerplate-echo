package transport

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *UserHandler) {
	userHandler := NewUserHandler()
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/test", userHandler.Test)

	}
}
