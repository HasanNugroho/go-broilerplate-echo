package auth

import (
	"net/http"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService IAuthService
	app         *app.Apps
}

func NewAuthHandler(us IAuthService, app *app.Apps) *AuthHandler {
	return &AuthHandler{
		authService: us,
		app:         app,
	}
}

// Register godoc
// @Summary      Register
// @Description  Register an user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  users.UserCreateModel  true  "User Data"
// @Success      200 {object}  shared.Response{data=users.UserCreateModel}
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /auth/register [post]
func (c *AuthHandler) Register(ctx echo.Context) error {
	var user users.UserCreateModel
	if err := ctx.Bind(&user); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Failed to process request", "Invalid data format")
		return err
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Validation failed", err.Error())
		return err
	}

	err := c.authService.Register(ctx, c.app, &user)
	if err != nil {
		utils.SendError(ctx, http.StatusBadRequest, err.Error(), nil)
		return err
	}

	utils.SendSuccess(ctx, http.StatusOK, "Register successful", nil)
	return nil
}

// Login godoc
// @Summary      Login
// @Description  Login an user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  AuthModel  true  "User Data"
// @Success      200 {object}  shared.Response{data=AuthResponse}
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /auth/login [post]
func (c *AuthHandler) Login(ctx echo.Context) error {
	var user AuthModel
	if err := ctx.Bind(&user); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Failed to process request", "Invalid data format")
		return err
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Validation failed", err.Error())
		return err
	}

	token, err := c.authService.Login(ctx, c.app, user.Email, user.Password)
	if err != nil {
		utils.SendError(ctx, http.StatusUnauthorized, "Login failed", err.Error())
		return err
	}

	utils.SendSuccess(ctx, http.StatusOK, "Login successful", token)
	return nil
}

// Logout godoc
// @Summary      Logout
// @Description  Logout an user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LogoutRequest true "Logout payload"
// @Success      200 {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /auth/logout [post]
// @Security ApiKeyAuth
func (c *AuthHandler) Logout(ctx echo.Context) error {
	err := c.authService.Logout(ctx, c.app)
	if err != nil {
		utils.SendError(ctx, http.StatusUnauthorized, "Logout failed", err.Error())
		return err
	}

	utils.SendSuccess(ctx, http.StatusOK, "Logout successful", nil)
	return nil
}

// Renew token godoc
// @Summary      Renew token
// @Description  Renew token an user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LogoutRequest true "Logout payload"
// @Success      200 {object}  shared.Response
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /auth/access-token [post]
// @Security ApiKeyAuth
func (c *AuthHandler) GenerateAccessToken(ctx echo.Context) error {
	token, err := c.authService.GenerateAccessToken(ctx, c.app)
	if err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Renew token failed", err.Error())
		return err
	}

	utils.SendSuccess(ctx, http.StatusOK, "Renew token successfully", token)
	return nil
}
