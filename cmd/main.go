package main

import (
	"log"
	"net/http"

	"github.com/HasanNugroho/starter-golang/bootstrap"
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Initialize configuration
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Apply CORS middleware
	r.Use(middleware.SetCORS())

	// Get config
	configPtr := config.GetConfig()
	if configPtr == nil {
		log.Fatal("Failed to get config: config is nil")
	}
	config := *configPtr

	// Set production mode
	if config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize RDBMS
	if config.Database.RDBMS.Activate {
		db, err := bootstrap.InitDB()
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		// Close database connection
		defer bootstrap.ShutdownDB(db)
	}

	// Initialize REDIS client
	if config.Database.REDIS.Activate {
		redis, err := bootstrap.InitRedis()

		if err != nil {
			log.Fatalf("Failed to initialize redis: %v", err)
		}

		// Close redis connection
		defer bootstrap.ShutdownRedis(redis)
	}

	// Initialize Rate Limit
	if config.Security.RateLimit != "" {
		limiter, err := bootstrap.InitRateLimiter(config.Security.RateLimit, config.Security.TrustedPlatform)
		if err != nil {
			log.Fatalf("Failed to initialize rate-limit: %v", err)
		}
		r.Use(middleware.RateLimit(limiter))
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Fatal(r.Run(config.Server.ServerHost + ":" + config.Server.ServerPort))
}
