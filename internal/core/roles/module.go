package roles

import (
	"fmt"

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
	fmt.Println("Role Module Initialized")
	return nil
}

func (a *RoleModule) Route(router *gin.RouterGroup, app *app.Apps) {
	route := router.Group("/v1/roles")
	{
		route.Use(middleware.AuthMiddleware(app))
		route.POST("", a.Handler.Create)
		route.GET("", a.Handler.FindAll)
		route.GET(":id", a.Handler.FindById)
		route.PUT(":id", a.Handler.Update)
		route.DELETE(":id", a.Handler.Delete)
	}
}
