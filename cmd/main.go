package main

import (
	"log"
	"net/http"

	"github.com/HasanNugroho/starter-golang/bootstrap"
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/core/transport"
	"github.com/HasanNugroho/starter-golang/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	if err := config.InitConfig(); err != nil {
		log.Fatalf("❌ Failed to initialize config: %v", err)
	}

	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatal("❌ Failed to get config: config is nil")
	}

	// Set production mode if applicable
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router with middlewares
	r := gin.Default()
	r.Use(middleware.SetCORS(), middleware.SecurityMiddleware())

	// Initialize RDBMS if enabled
	if cfg.Database.RDBMS.Activate {
		db, err := bootstrap.InitDB(&cfg.Database.RDBMS)
		if err != nil {
			log.Fatalf("❌ Failed to initialize database: %v", err)
		}
		defer bootstrap.ShutdownDB(db) // Ensure database is closed on exit
	}

	// Initialize Redis if enabled
	if cfg.Database.REDIS.Activate {
		var err error
		redisClient, err := bootstrap.InitRedis()
		if err != nil {
			log.Fatalf("❌ Failed to initialize Redis: %v", err)
		}
		cfg.Database.REDIS.Client = redisClient
		defer bootstrap.ShutdownRedis(redisClient)
	}

	// Initialize Rate Limiter if enabled
	if cfg.Security.RateLimit != "" {
		limiter, err := bootstrap.InitRateLimiter(cfg, cfg.Security.RateLimit, cfg.Security.TrustedPlatform)
		if err != nil {
			log.Fatalf("❌ Failed to initialize rate limiter: %v", err)
		}
		r.Use(middleware.RateLimit(limiter))
	}

	// Define test route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	transport.RegisterRoutes(r, &transport.UserHandler{})

	// Start server
	serverAddr := cfg.Server.ServerHost + ":" + cfg.Server.ServerPort
	log.Printf("🚀 Server running at %s", serverAddr)
	log.Fatal(r.Run(serverAddr))
}
