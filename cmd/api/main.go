package main

import (
	"net"

	"github.com/HasanNugroho/starter-golang/cmd/docs"
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

	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Initialize Rate Limiter if enabled
	if cfg.Security.RateLimit != "" {
		limiter, err := config.InitRateLimiter(cfg, cfg.Security.RateLimit, cfg.Security.TrustedPlatform)
		if err != nil {
			config.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		r.Use(middleware.RateLimit(limiter))
	}

	route, _ := InitializeApp(r, &cfg.Database)
	route.SetupRoutes()

	// Start server
	serverAddr := net.JoinHostPort(cfg.Server.ServerHost, cfg.Server.ServerPort)
	// config.Logger.Info().Msg("🚀 Server running at " + serverAddr)
	err = r.Run(serverAddr)
	if err != nil {
		config.Logger.Fatal().Msg(err.Error())
	}

}
