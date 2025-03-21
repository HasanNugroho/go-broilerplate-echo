package main

import (
	"fmt"

	"github.com/HasanNugroho/starter-golang/docs"
	config "github.com/HasanNugroho/starter-golang/internal/configs"
	"github.com/HasanNugroho/starter-golang/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	ProductionEnv = "production"
)

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
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal().Msg("❌ Failed to initialize config: " + err.Error())
	}

	// Set production mode if applicable
	if cfg.AppEnv == ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Logger
	config.InitLogger(cfg)

	// Initialize RDBMS if enabled
	if cfg.Database.Enabled {
		db, err := config.InitDB(&cfg.Database)
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		defer config.ShutdownDB(db) // Ensure database is closed on exit
	}

	// Initialize Redis if enabled
	if cfg.Redis.Enabled {
		var err error
		redisClient, err := config.InitRedis(&cfg.Redis)
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		// config.Database.Redis.Client = redisClient
		defer config.ShutdownRedis(redisClient)
	}

	r := config.NewGin(cfg)
	r.Use(middleware.SetCORS(cfg), middleware.SecurityMiddleware(cfg))

	loadSwagger(r, cfg)

	// Initialize Rate Limiter if enabled
	if cfg.Security.RateLimit != "" {
		limiter, err := config.InitRateLimiter(cfg, cfg.Security.RateLimit, cfg.Security.TrustedPlatform)
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		r.Use(middleware.RateLimit(limiter))
	}

	route, err := InitializeRoute(r, &cfg.Database)

	if err != nil {
		config.Logger.Fatal().Msg("❌ Failed to initialize routes: " + err.Error())
		panic(1)
	}
	route.SetupRoutes()

	for _, route := range r.Routes() {
		fmt.Println("Registered Route:", route.Method, route.Path)
	}

	err = r.Run(fmt.Sprintf(":%s", cfg.Server.ServerPort))
	if err != nil {
		config.Logger.Fatal().Msg(err.Error())
	}

}

func loadSwagger(r *gin.Engine, cfg *config.Configuration) {
	docs.SwaggerInfo.Title = cfg.APP_NAME
	docs.SwaggerInfo.Version = cfg.Version
	docs.SwaggerInfo.Description = cfg.APP_NAME
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.ServerHost, cfg.Server.ServerPort)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
