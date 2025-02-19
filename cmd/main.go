package main

import (
	"log"
	"net/http"

	"github.com/HasanNugroho/starter-golang/pkg/config"
	"github.com/HasanNugroho/starter-golang/pkg/database"
	"github.com/HasanNugroho/starter-golang/pkg/lib"
	"github.com/HasanNugroho/starter-golang/pkg/middleware"
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

	// limiterInstance, err := config.(
	// 	configure.Security.RateLimit,
	// 	trustedPlatform,
	// )

	configPtr := config.GetConfig()
	if configPtr == nil {
		log.Fatal("Failed to get config: config is nil")
	}
	config := *configPtr

	// Initialize RDBMS
	if config.Database.RDBMS.Activate {
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		// Close database connection
		defer database.ShutdownDB(db)
	}

	// Initialize REDIS client
	if config.Database.REDIS.Activate {
		redis, err := database.InitRedis()

		if err != nil {
			log.Fatalf("Failed to initialize redis: %v", err)
		}

		// Close redis connection
		defer database.ShutdownRedis(redis)
	}

	// Initialize Rate Limit
	if config.Security.RateLimit != "" {
		limiter, err := lib.InitRateLimiter(config.Security.RateLimit, config.Security.TrustedPlatform)
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
