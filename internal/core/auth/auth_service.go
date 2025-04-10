package auth

import (
	"fmt"
	"net/http"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/entities"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthService struct {
	repo users.IUserRepository
}

func NewAuthService(repo users.IUserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) Login(ctx echo.Context, app *app.Apps, email string, password string) (AuthResponse, error) {
	existingUser, err := a.repo.FindByEmail(ctx, email)
	if err != nil || existingUser.Email == "" {
		return AuthResponse{}, fmt.Errorf("Incorrect email or password %w", err)
	}

	if !utils.VerifyPassword(existingUser.Password, []byte(password)) {
		return AuthResponse{}, fmt.Errorf("Incorrect email or password")
	}

	var allPermissions []string
	for _, role := range existingUser.RolesData {
		allPermissions = append(allPermissions, role.Permissions...)
	}

	payload := map[string]interface{}{
		"id":         existingUser.ID,
		"email":      existingUser.Email,
		"name":       existingUser.Name,
		"created_at": existingUser.CreatedAt,
		"permission": allPermissions,
		"roles":      existingUser.RolesData,
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

func (a *AuthService) Register(ctx echo.Context, app *app.Apps, user *users.UserCreateModel) error {
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

func (a *AuthService) Logout(ctx echo.Context, app *app.Apps) error {
	tokenString := ctx.Request().Header.Get("Authorization")
	if tokenString == "" {
		return fmt.Errorf("Token is required")
	}

	var req LogoutRequest
	if err := ctx.Bind(&req); err != nil || req.RefreshToken == "" {
		return fmt.Errorf("refresh token is required")
	}

	err := utils.RevokeToken(app, tokenString, req.RefreshToken)
	if err != nil {
		return fmt.Errorf("Failed to revoke token")
	}

	return nil
}

func (a *AuthService) GenerateAccessToken(ctx echo.Context, app *app.Apps) (AuthResponse, error) {
	tokenString := ctx.Request().Header.Get("Authorization")
	token, err := utils.ValidateToken(app, tokenString)
	if err != nil || token == nil || !token.Valid {
		return AuthResponse{}, fmt.Errorf("invalid or expired access token: %w", err)
	}

	// Get and validate Authorization header (old access token)
	claimsRaw := ctx.Get("claims")
	if claimsRaw == nil {
		return AuthResponse{}, fmt.Errorf("No claims found")
	}

	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		return AuthResponse{}, fmt.Errorf("Invalid claims format")
	}

	data, ok := claims["data"].(map[string]interface{})
	if !ok {
		utils.SendError(ctx, http.StatusForbidden, "Invalid data in claims", nil)
	}

	userID, ok := data["id"].(string)
	if !ok {
		return AuthResponse{}, fmt.Errorf("invalid or missing user ID in token")
	}

	existingUser, err := a.repo.FindById(ctx, userID)

	app.Log.Info().Msgf("User ID from token: %s", existingUser)
	if err != nil || existingUser.Email == "" {
		return AuthResponse{}, fmt.Errorf("user not found")
	}

	// Aggregate permissions
	var allPermissions []string
	for _, role := range existingUser.RolesData {
		allPermissions = append(allPermissions, role.Permissions...)
	}

	// Parse refresh token from request body
	var req LogoutRequest
	if err := ctx.Bind(&req); err != nil || req.RefreshToken == "" {
		return AuthResponse{}, fmt.Errorf("refresh token is required")
	}

	newPayload := map[string]interface{}{
		"id":          userID,
		"email":       existingUser.Email,
		"name":        existingUser.Name,
		"permissions": allPermissions,
	}

	// Generate new access token
	newAccessToken, err := utils.RefreshAccessToken(app, req.RefreshToken, newPayload)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to refresh token: %w", err)
	}

	return AuthResponse{
		Token:        newAccessToken,
		RefreshToken: req.RefreshToken,
		Data:         newPayload,
	}, nil
}
