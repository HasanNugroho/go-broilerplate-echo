package users

import (
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
	app.Log.Info().Msg("User Module Initialized")

	permission := []string{
		"users:create",
		"users:read",
		"users:update",
		"users:delete",
		"manage:system",
	}

	// Merge permission
	app.Config.ModulePermissions = append(app.Config.ModulePermissions, permission...)

	return nil
}

func (u *UserModule) Route(router *gin.RouterGroup, app *app.Apps) {
	userRoutes := router.Group("/v1/users")
	userRoutes.Use(middleware.AuthMiddleware(app))
	{
		userRoutes.POST("/", middleware.CheckAccess([]string{"users:create"}), u.Handler.Create)
		userRoutes.GET("/", middleware.CheckAccess([]string{"users:read"}), u.Handler.FindAll)
		userRoutes.GET("/:id", middleware.CheckAccess([]string{"users:read"}), u.Handler.FindById)
		userRoutes.PUT("/:id", middleware.CheckAccess([]string{"users:update"}), u.Handler.Update)
		userRoutes.DELETE("/:id", middleware.CheckAccess([]string{"users:delete"}), u.Handler.Delete)
	}
}
