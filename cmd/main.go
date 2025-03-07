package main

import (
	"net/http"

	"github.com/HasanNugroho/starter-golang/bootstrap"
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	ProductionEnv = "production"
)

func main() {
	// Initialize configuration
	config, err := config.InitConfig()
	if err != nil {
		log.Fatal().Msg("❌ Failed to initialize config: " + err.Error())
	}

	// Set production mode if applicable
	if config.AppEnv == ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Logger
	bootstrap.InitLogger(config)

	// Initialize Gin router with middlewares
	r := gin.Default()
	r.Use(middleware.SetCORS(config), middleware.SecurityMiddleware(config))

	// Initialize RDBMS if enabled
	if config.Database.RDBMS.Enabled {
		db, err := bootstrap.InitDB(&config.Database.RDBMS)
		if err != nil {
			bootstrap.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		defer bootstrap.ShutdownDB(db) // Ensure database is closed on exit
	}

	// Initialize Redis if enabled
	if config.Database.Redis.Enabled {
		var err error
		redisClient, err := bootstrap.InitRedis(&config.Database.Redis)
		if err != nil {
			bootstrap.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		// config.Database.Redis.Client = redisClient
		defer bootstrap.ShutdownRedis(redisClient)
	}

	// Initialize Rate Limiter if enabled
	if config.Security.RateLimit != "" {
		limiter, err := bootstrap.InitRateLimiter(config, config.Security.RateLimit, config.Security.TrustedPlatform)
		if err != nil {
			bootstrap.Logger.Fatal().Msg(err.Error())
			panic(1)
		}
		r.Use(middleware.RateLimit(limiter))
	}

	// Define test route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Start server
	serverAddr := config.Server.ServerHost + ":" + config.Server.ServerPort
	// bootstrap.Logger.Info().Msg("🚀 Server running at " + serverAddr)
	err = r.Run(serverAddr)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
