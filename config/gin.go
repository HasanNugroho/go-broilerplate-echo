package config

import (
	"github.com/gin-gonic/gin"
)

func NewGin() *gin.Engine {
	// Initialize Gin router
	r := gin.Default()
	return r
}
