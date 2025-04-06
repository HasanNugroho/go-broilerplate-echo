package main

import (
	"fmt"

	"github.com/HasanNugroho/starter-golang/cmd/docs"
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	ProductionEnv = "production"
)

// @title           Example Rest API
// @version         1.0
// @description     This is a sample server celler server.

// @host      localhost:7000
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Initialize configuration and other components
	router := gin.Default()
	apps := internal.AppsInit(router)

	loadSwagger(apps.Router, apps.Config)

	err := apps.Router.Run(fmt.Sprintf(":%s", apps.Config.Server.ServerPort))
	if err != nil {
		config.Logger.Fatal().Msg(err.Error())
	}

}

func loadSwagger(r *gin.Engine, appConfig *config.Config) {
	docs.SwaggerInfo.Title = appConfig.AppName
	docs.SwaggerInfo.Version = appConfig.Version
	docs.SwaggerInfo.Description = appConfig.AppName
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", appConfig.Server.ServerHost, appConfig.Server.ServerPort)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
