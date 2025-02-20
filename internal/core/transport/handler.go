package transport

import (
	"context"
	"log"

	"github.com/HasanNugroho/starter-golang/config"
	utils "github.com/HasanNugroho/starter-golang/pkg/utlis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// UserHandler handles user-related operations
type UserHandler struct {
	redisClient *redis.Client
}

// NewUserHandler initializes the UserHandler
func NewUserHandler() *UserHandler {
	cfg := config.GetConfig()
	if cfg.Database.REDIS.Client == nil {
		log.Fatal("❌ Redis client is not initialized")
	}

	return &UserHandler{
		redisClient: cfg.Database.REDIS.Client,
	}
}

var ctx = context.Background()

// Test retrieves and sets a value in Redis
func (h *UserHandler) Test(c *gin.Context) {
	const key = "test"

	// Try getting the value from Redis
	val, err := h.redisClient.Get(ctx, key).Result()
	if err == redis.Nil { // Key not found
		if err := h.redisClient.Set(ctx, key, "testing", 0).Err(); err != nil {
			utils.SendError(c, 500, "Failed to set value in Redis", err.Error())
			return
		}
		val, _ = h.redisClient.Get(ctx, key).Result()
	} else if err != nil { // Other Redis errors
		utils.SendError(c, 500, "Redis error", err.Error())
		return
	}

	// Send success response
	utils.SendSuccess(c, "Success", val)
}
