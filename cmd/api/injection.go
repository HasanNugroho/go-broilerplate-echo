//go:build wireinject
// +build wireinject

package main

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal"
	"github.com/HasanNugroho/starter-golang/internal/core/auth"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	users.NewUserRepository,
	wire.Bind(new(users.IUserRepository), new(*users.UserRepository)),
	users.NewUserService,
	wire.Bind(new(users.IUserService), new(*users.UserService)),
	users.NewUserHandler,
)

var authSet = wire.NewSet(
	auth.NewAuthService,
	wire.Bind(new(auth.IAuthService), new(*auth.AuthService)),
	auth.NewAuthHandler,
)

func InitializeRoute(r *gin.Engine, cfg *config.DatabaseConfig) (*internal.RouteConfig, error) {
	wire.Build(
		userSet,
		authSet,
		internal.NewRouter,
	)

	return nil, nil
}
