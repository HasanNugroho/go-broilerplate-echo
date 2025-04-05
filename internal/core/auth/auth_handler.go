package auth

import (
	"net/http"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService IAuthService
	config      *config.Config
}

func NewAuthHandler(us IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: us,
		config:      config.GetConfig(),
	}
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
func (c *AuthHandler) Login(ctx *gin.Context) {
	var user AuthModel
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Failed to process request", "Invalid data format")
		return
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	token, err := c.authService.Login(ctx, c.config, user.Email, user.Password)
	if err != nil {
		utils.SendError(ctx, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	utils.SendSuccess(ctx, http.StatusOK, "Login successful", token)
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
func (c *AuthHandler) Logout(ctx *gin.Context) {
	err := c.authService.Logout(ctx, c.config)
	if err != nil {
		utils.SendError(ctx, http.StatusUnauthorized, "Logout failed", err.Error())
		return
	}

	utils.SendSuccess(ctx, http.StatusOK, "Logout successful", nil)
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
func (c *AuthHandler) GenerateAccessToken(ctx *gin.Context) {
	token, err := c.authService.GenerateAccessToken(ctx, c.config)
	if err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Renew token failed", err.Error())
		return
	}

	utils.SendSuccess(ctx, http.StatusOK, "Renew token successfully", token)
}
