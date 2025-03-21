//go:build wireinject
// +build wireinject

package main

import (
	"github.com/HasanNugroho/starter-golang/internal"
	"github.com/HasanNugroho/starter-golang/internal/configs"
	"github.com/HasanNugroho/starter-golang/internal/users"
	"github.com/HasanNugroho/starter-golang/internal/users/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	wire.Bind(new(repository.IUserRepository), new(*repository.UserRepository)),
	users.NewUserService, // Pastikan `UserService` di dalam package `users`
	wire.Bind(new(users.IUserService), new(*users.UserService)), // Bind ke IUserService
	users.NewUserHandler,
)

func InitializeApp(r *gin.Engine, cfg *configs.RDBMSConfig) (*internal.RouteConfig, error) {
	wire.Build(
		userSet,
		internal.NewRouter,
	)

	return nil, nil
}
