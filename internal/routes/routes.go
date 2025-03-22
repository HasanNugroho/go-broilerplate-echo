package routes

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/users"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App         *gin.Engine
	UserHandler *users.UserHandler
	Logger      *config.LoggerConfig
}

func NewRouter(app *gin.Engine, userHandler *users.UserHandler) *RouteConfig {
	return &RouteConfig{
		App:         app,
		UserHandler: userHandler,
	}
}

func (r *RouteConfig) SetupRoutes() {
	v1 := r.App.Group("/api/v1")

	users.RegisterUserRoutes(v1, r.UserHandler)
}
