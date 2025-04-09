package roles

import (
	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/shared/middleware"
	"github.com/gin-gonic/gin"
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

func (a *RoleModule) Route(router *gin.RouterGroup, app *app.Apps) {
	route := router.Group("/v1/roles")
	{
		route.Use(middleware.AuthMiddleware(app))
		route.POST("", middleware.CheckAccess([]string{"roles:create"}), a.Handler.Create)
		route.GET("", middleware.CheckAccess([]string{"roles:read", "roles:assign", "roles:unassign"}), a.Handler.FindAll)
		route.GET(":id", middleware.CheckAccess([]string{"roles:read", "roles:assign", "roles:unassign"}), a.Handler.FindById)
		route.PUT(":id", middleware.CheckAccess([]string{"roles:update"}), a.Handler.Update)
		route.DELETE(":id", middleware.CheckAccess([]string{"roles:delete"}), a.Handler.Delete)
		route.POST("/assign", middleware.CheckAccess([]string{"roles:assign"}), a.Handler.AssignUser)
		route.POST("/unassign", middleware.CheckAccess([]string{"roles:unassign"}), a.Handler.AssignUser)
	}
}
