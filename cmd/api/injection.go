//go:build wireinject
// +build wireinject

package main

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/auth"
	"github.com/HasanNugroho/starter-golang/internal/routes"
	"github.com/HasanNugroho/starter-golang/internal/users"
	"github.com/HasanNugroho/starter-golang/internal/users/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	wire.Bind(new(repository.IUserRepository), new(*repository.UserRepository)),
	users.NewUserService,
	wire.Bind(new(users.IUserService), new(*users.UserService)),
	users.NewUserHandler,
)

var authSet = wire.NewSet(
	auth.NewAuthService,
	wire.Bind(new(auth.IAuthService), new(*auth.AuthService)),
	auth.NewAuthHandler,
)

func InitializeRoute(r *gin.Engine, cfg *config.DatabaseConfig) (*routes.RouteConfig, error) {
	wire.Build(
		userSet,
		authSet,
		routes.NewRouter,
	)

	return nil, nil
}
