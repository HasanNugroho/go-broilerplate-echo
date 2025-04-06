package users

import (
	"fmt"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

type UserModule struct {
	Handler *UserHandler
}

func NewUserModule(apps *app.Apps) *UserModule {
	userRepository := NewUserRepository(apps)
	userService := NewUserService(userRepository)
	userHandler := NewUserHandler(userService)
	return &UserModule{
		Handler: userHandler,
	}
}

func (u *UserModule) Register(app *app.Apps) error {
	fmt.Println("User Module Initialized")
	return nil
}

func (u *UserModule) Route(router *gin.RouterGroup, app *app.Apps) {
	userRoutes := router.Group("/v1/users")
	userRoutes.Use(middleware.AuthMiddleware(app))
	{
		userRoutes.POST("/", u.Handler.Create)
		userRoutes.GET("/", u.Handler.FindAll)
		userRoutes.GET("/:id", u.Handler.FindById)
		userRoutes.PUT("/:id", u.Handler.Update)
		userRoutes.DELETE("/:id", u.Handler.Delete)
	}
}
