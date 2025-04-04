package auth

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/core/auth/model"
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
// @Param        user  body  model.AuthModel  true  "User Data"
// @Success      200 {object}  shared.Response{data=model.AuthResponse}
// @Failure      400  {object}  shared.Response
// @Failure      404  {object}  shared.Response
// @Failure      500  {object}  shared.Response
// @Router       /auth/login [post]
func (c *AuthHandler) Login(ctx *gin.Context) {
	var user model.AuthModel
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.SendError(ctx, 400, "Failed to process request", "Invalid data format")
		return
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		utils.SendError(ctx, 400, "Validation failed", err.Error())
		return
	}

	token, err := c.authService.Login(ctx, c.config, user.Email, user.Password)
	if err != nil {
		utils.SendError(ctx, 401, "Login failed", err.Error())
		return
	}

	utils.SendSuccess(ctx, 200, "Login successful", token)
}
