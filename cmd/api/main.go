package main

import (
	"fmt"

	"github.com/HasanNugroho/starter-golang/cmd/docs"
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title       Starter Golang API
// @version     1.0
// @description This is a sample server.

// @contact.name   API Support
// @contact.email  support@example.com

// @host      localhost:7000
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Setup Echo and App
	router := echo.New()
	apps := internal.AppsInit(router)

	// Middleware
	// router.Use(middleware.ErrorHandler(apps))

	// Swagger Setup
	loadSwagger(router, apps.Config)

	// Start Server
	err := router.Start(fmt.Sprintf(":%s", apps.Config.Server.ServerPort))
	if err != nil {
		config.Logger.Fatal().Msg(err.Error())
	}
}

func loadSwagger(r *echo.Echo, cfg *config.Config) {
	docs.SwaggerInfo.Title = cfg.AppName
	docs.SwaggerInfo.Version = cfg.Version
	docs.SwaggerInfo.Description = cfg.AppName
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.ServerHost, cfg.Server.ServerPort)
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Swagger endpoint
	r.GET("/swagger/*", echoSwagger.WrapHandler)
}
