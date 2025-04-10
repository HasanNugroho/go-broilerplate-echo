package users

import (
	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/shared/middleware"
	"github.com/labstack/echo/v4"
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

func (u *UserModule) Route(router *echo.Group, app *app.Apps) {
	userRoutes := router.Group("/v1/users")
	userRoutes.Use(middleware.AuthMiddleware(app))
	{
		userRoutes.POST("", u.Handler.Create, middleware.CheckAccess([]string{"users:create"}))
		userRoutes.GET("/", u.Handler.FindAll, middleware.CheckAccess([]string{"users:read"}))
		userRoutes.GET("/:id", u.Handler.FindById, middleware.CheckAccess([]string{"users:read"}))
		userRoutes.PUT("/:id", u.Handler.Update, middleware.CheckAccess([]string{"users:update"}))
		userRoutes.DELETE("/:id", u.Handler.Delete, middleware.CheckAccess([]string{"users:delete"}))

	}
}
