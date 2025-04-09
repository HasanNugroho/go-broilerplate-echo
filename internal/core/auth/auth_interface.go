package auth

import (
	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Login(ctx *gin.Context, app *app.Apps, email string, password string) (AuthResponse, error)
	Register(ctx *gin.Context, app *app.Apps, user *users.UserCreateModel) error
	Logout(ctx *gin.Context, app *app.Apps) error
	GenerateAccessToken(ctx *gin.Context, app *app.Apps) (AuthResponse, error)
}
