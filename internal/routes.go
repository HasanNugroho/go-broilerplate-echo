package internal

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/core/auth"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App         *gin.Engine
	UserHandler *users.UserHandler
	AuthHandler *auth.AuthHandler
	Logger      *config.LoggerConfig
	Config      *config.Config
}

func NewRouter(app *gin.Engine, config *config.Config, userHandler *users.UserHandler, authHandler *auth.AuthHandler) *RouteConfig {
	return &RouteConfig{
		App:         app,
		UserHandler: userHandler,
		AuthHandler: authHandler,
		Config:      config,
	}
}

func (r *RouteConfig) SetupRoutes() {
	v1 := r.App.Group("/api/v1")

	auth.RegisterAuthRoutes(v1, r.Config, r.AuthHandler)

	users.RegisterUserRoutes(v1, r.Config, r.UserHandler)
}
