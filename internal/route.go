package internal

import (
	"github.com/HasanNugroho/starter-golang/internal/configs"
	"github.com/HasanNugroho/starter-golang/internal/users"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App         *gin.Engine
	UserHandler *users.UserHandler
	Logger      *configs.LoggerConfig
}

func NewRouter(app *gin.Engine, userHandler *users.UserHandler) *RouteConfig {
	return &RouteConfig{
		App:         app,
		UserHandler: userHandler,
	}
}

func (r *RouteConfig) SetupRoutes() {
	r.App.POST("/user", r.UserHandler.Create)
}
