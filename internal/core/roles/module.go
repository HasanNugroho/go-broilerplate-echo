package roles

import (
	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type RoleModule struct {
	Handler *RoleHandler
}

func NewRoleModule(app *app.Apps) *RoleModule {
	roleRepository := NewRoleRepository(app)
	roleService := NewRoleService(app, roleRepository)
	roleHandler := NewRoleHandler(roleService)
	return &RoleModule{
		Handler: roleHandler,
	}
}

func (u *RoleModule) Register(app *app.Apps) error {
	app.Log.Info().Msg("Role Module Initialized")

	permission := []string{
		"roles:unassign",
		"roles:read",
		"roles:update",
		"roles:delete",
		"roles:assign",
		"roles:unassign",
	}

	// Merge permission
	app.Config.ModulePermissions = append(app.Config.ModulePermissions, permission...)

	return nil
}

func (a *RoleModule) Route(router *echo.Group, app *app.Apps) {
	route := router.Group("/v1/roles")
	{
		route.Use(middleware.AuthMiddleware(app))
		route.POST("", a.Handler.Create, middleware.CheckAccess([]string{"roles:create"}))
		route.GET("", a.Handler.FindAll, middleware.CheckAccess([]string{"roles:read", "roles:assign", "roles:unassign"}))
		route.GET("/:id", a.Handler.FindById, middleware.CheckAccess([]string{"roles:read", "roles:assign", "roles:unassign"}))
		route.PUT("/:id", a.Handler.Update, middleware.CheckAccess([]string{"roles:update"}))
		route.DELETE("/:id", a.Handler.Delete, middleware.CheckAccess([]string{"roles:delete"}))
		route.POST("/assign", a.Handler.AssignUser, middleware.CheckAccess([]string{"roles:assign"}))
		route.POST("/unassign", a.Handler.UnAssignUser, middleware.CheckAccess([]string{"roles:unassign"}))

	}
}
