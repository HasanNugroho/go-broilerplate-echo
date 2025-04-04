package auth

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Login(ctx *gin.Context, config *config.Config, email string, password string) (AuthResponse, error)
	Logout(ctx *gin.Context, config *config.Config) error
	GenerateAccessToken(ctx *gin.Context, config *config.Config) (AuthResponse, error)
	Register(ctx *gin.Context, email string, password string) (interface{}, error)
	ResetPassword(ctx *gin.Context, email string, password string) error
}
