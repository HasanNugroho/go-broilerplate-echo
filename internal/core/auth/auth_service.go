package auth

import (
	"encoding/json"
	"fmt"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
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

func (a *AuthService) Login(ctx *gin.Context, app *app.Apps, email string, password string) (AuthResponse, error) {
	existingUser, err := a.repo.FindByEmail(ctx, email)
	if err != nil || existingUser.Email == "" {
		return AuthResponse{}, fmt.Errorf("Incorrect email or password")
	}

	if !utils.VerifyPassword(existingUser.Password, []byte(password)) {
		return AuthResponse{}, fmt.Errorf("Incorrect email or password")
	}

	var allPermissions []string
	for _, role := range existingUser.Roles {

		var perms []string
		err := json.Unmarshal([]byte(role.Permissions), &perms)
		if err != nil {
			panic(err)
		}

		allPermissions = append(allPermissions, perms...)
	}

	payload := map[string]interface{}{
		"id":         (existingUser.ID).String(),
		"email":      existingUser.Email,
		"name":       existingUser.Name,
		"created_at": existingUser.CreatedAt,
		"permission": allPermissions,
		"roles":      existingUser.Roles,
	}

	accessToken, refreshToken, err := utils.GenerateAuthToken(app, payload)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("Error creating token: %s", err.Error())
	}

	return AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		Data:         payload,
	}, nil
}

func (a *AuthService) Register(ctx *gin.Context, app *app.Apps, user *users.UserCreateModel) error {
	existingUser, err := a.repo.FindByEmail(ctx, user.Email)
	if err != nil || existingUser.Email != "" {
		return fmt.Errorf("Incorrect email or password")
	}

	password, err := utils.HashPassword([]byte(user.Password))
	if err != nil {
		return err
	}

	payload := entities.User{
		Email:    user.Email,
		Name:     user.Name,
		Password: password,
	}

	if err = a.repo.Create(ctx, &payload); err != nil {
		return err
	}
	return nil
}

func (a *AuthService) Logout(ctx *gin.Context, app *app.Apps) error {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		return fmt.Errorf("Token is required")
	}

	var req LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		return fmt.Errorf("refresh token is required")
	}

	err := utils.RevokeToken(app, tokenString, req.RefreshToken)
	if err != nil {
		return fmt.Errorf("Failed to revoke token")
	}

	return nil
}

func (a *AuthService) GenerateAccessToken(ctx *gin.Context, app *app.Apps) (AuthResponse, error) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		return AuthResponse{}, fmt.Errorf("token is required")
	}

	// Parse the token
	token, err := utils.ValidateToken(app, tokenString)
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
	newAccessToken, err := utils.RefreshAccessToken(app, req.RefreshToken)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to refresh token: %v", err)
	}
	return AuthResponse{
		Token:        newAccessToken,
		RefreshToken: req.RefreshToken,
		Data:         data,
	}, nil
}
