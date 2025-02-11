package main

import (
	"fmt"
	"log"

	"github.com/HasanNugroho/starter-golang/pkg/config"
	"github.com/HasanNugroho/starter-golang/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	configPtr := config.GetConfig()
	if configPtr == nil {
		log.Fatal("Failed to get config: config is nil")
	}
	config := *configPtr

	if config.Database.RDBMS.Activate {
		db, err := database.InitDB()
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		fmt.Println("Database connected successfully:", db)
	}

	if config.Database.REDIS.Activate {
		// Initialize REDIS client
		err := database.InitRedis()

		if err != nil {
			log.Fatalf("Failed to initialize redis: %v", err)
		}
	}

	log.Fatal(r.Run(config.Server.ServerHost + ":" + config.Server.ServerPort))
}
