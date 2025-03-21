package configs

import (
	"github.com/gin-gonic/gin"
)

func NewGin(cfg *Configuration) *gin.Engine {
	// Initialize Gin router
	r := gin.Default()
	return r
}
