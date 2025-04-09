package auth

import (
	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/labstack/echo/v4"
)

type IAuthService interface {
	Login(ctx echo.Context, app *app.Apps, email string, password string) (AuthResponse, error)
	Register(ctx echo.Context, app *app.Apps, user *users.UserCreateModel) error
	Logout(ctx echo.Context, app *app.Apps) error
	GenerateAccessToken(ctx echo.Context, app *app.Apps) (AuthResponse, error)
}
