package main

import (
	"fmt"

	"github.com/HasanNugroho/starter-golang/cmd/docs"
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	ProductionEnv = "production"
)

var appConfig *config.Config

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:7000
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Initialize configuration
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msg("❌ Failed to initialize config: " + err.Error())
	}

	// Set production mode if applicable
	if appConfig.AppEnv == ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Logger
	config.InitLogger(appConfig)

	// Initialize RDBMS if enabled
	if appConfig.DB.Enabled {
		db, err := appConfig.DB.InitDB()
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		defer config.ShutdownDB(db)
	}

	// Initialize Redis if enabled
	if appConfig.Redis.Enabled {
		redisClient, err := appConfig.Redis.InitRedis()
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		defer config.ShutdownRedis(redisClient)
	}

	// Initialize Elastic if enabled
	if appConfig.Search.Enabled {
		err := appConfig.Search.SearchInit()
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
	}

	r := config.NewGin(appConfig)
	r.Use(middleware.SetCORS(appConfig), middleware.SecurityMiddleware(appConfig))

	loadSwagger(r, appConfig)

	// Initialize Rate Limiter if enabled
	if appConfig.Security.RateLimit != "" {
		limiter, err := config.InitRateLimiter(appConfig, appConfig.Security.RateLimit, appConfig.Security.TrustedPlatform)
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		r.Use(middleware.RateLimit(limiter))
	}

	route, err := InitializeRoute(r, &appConfig.DB)

	if err != nil {
		config.Logger.Fatal().Msg("❌ Failed to initialize routes: " + err.Error())
		panic(1)
	}
	route.SetupRoutes()

	for _, route := range r.Routes() {
		fmt.Println("Registered Route:", route.Method, route.Path)
	}

	err = r.Run(fmt.Sprintf(":%s", appConfig.Server.ServerPort))
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
