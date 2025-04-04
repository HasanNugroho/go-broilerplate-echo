package auth

import (
	"fmt"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	repo users.IUserRepository
}

func NewAuthService(repo users.IUserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) Login(ctx *gin.Context, config *config.Config, email string, password string) (AuthResponse, error) {
	existingUser, err := a.repo.FindByEmail(ctx, email)
	if err != nil || existingUser.Email == "" {
		return AuthResponse{}, fmt.Errorf("Incorrect email or password")
	}

	if !utils.VerifyPassword(existingUser.Password, []byte(password)) {
		return AuthResponse{}, fmt.Errorf("Incorrect email or password")
	}

	payload := users.UserModelResponse{
		ID:        (existingUser.ID).String(),
		Email:     existingUser.Email,
		Name:      existingUser.Name,
		CreatedAt: existingUser.CreatedAt,
	}

	accessToken, refreshToken, err := utils.GenerateAuthToken(config, payload)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("Error creating token: %s", err.Error())
	}

	return AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		Data:         payload,
	}, nil
}

func (a *AuthService) Logout(ctx *gin.Context, config *config.Config) error {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		return fmt.Errorf("Token is required")
	}

	var req LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		return fmt.Errorf("refresh token is required")
	}

	err := utils.RevokeToken(config, tokenString, req.RefreshToken)
	if err != nil {
		return fmt.Errorf("Failed to revoke token")
	}

	return nil
}

func (a *AuthService) GenerateAccessToken(ctx *gin.Context, config *config.Config) (AuthResponse, error) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		return AuthResponse{}, fmt.Errorf("token is required")
	}

	// Parse the token
	token, err := utils.ValidateToken(config, tokenString)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("invalid refresh token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AuthResponse{}, fmt.Errorf("cannot parse claims")
	}
	data, ok := claims["data"].(map[string]interface{})
	if !ok {
		return AuthResponse{}, fmt.Errorf("data claim not found or invalid format")
	}

	var req LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		return AuthResponse{}, fmt.Errorf("refresh token is required")
	}
	newAccessToken, err := utils.RefreshAccessToken(config, req.RefreshToken)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to refresh token: %v", err)
	}
	return AuthResponse{
		Token:        newAccessToken,
		RefreshToken: req.RefreshToken,
		Data:         data,
	}, nil
}

func (a *AuthService) Register(ctx *gin.Context, email string, password string) (interface{}, error) {
	panic("not implemented") // TODO: Implement
}

func (a *AuthService) ResetPassword(ctx *gin.Context, email string, password string) error {
	panic("not implemented") // TODO: Implement
}
