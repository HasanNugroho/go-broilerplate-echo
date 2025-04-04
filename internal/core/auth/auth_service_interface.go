package auth

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/core/auth/model"
	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Login(ctx *gin.Context, config *config.Config, email string, password string) (model.AuthResponse, error)
	Register(ctx *gin.Context, email string, password string) (interface{}, error)
	ResetPassword(ctx *gin.Context, email string, password string) error
}
