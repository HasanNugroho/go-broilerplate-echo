package auth

import (
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
		return AuthResponse{}, utils.NewBadRequest("Incorrect email or password")
	}

	if !utils.VerifyPassword(existingUser.Password, []byte(password)) {
		return AuthResponse{}, utils.NewBadRequest("Incorrect email or password")
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
		return AuthResponse{}, utils.NewInternal(err.Error())
	}

	return AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		Data:         payload,
	}, nil
}

func (a *AuthService) Register(ctx echo.Context, app *app.Apps, user *users.UserCreateModel) error {
	_, err := a.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
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
		return utils.NewBadRequest("Token is required")
	}

	var req LogoutRequest
	if err := ctx.Bind(&req); err != nil || req.RefreshToken == "" {
		return utils.NewBadRequest("refresh token is required")
	}

	err := utils.RevokeToken(app, tokenString, req.RefreshToken)
	if err != nil {
		return utils.NewInternal("Failed to revoke token")
	}

	return nil
}

func (a *AuthService) GenerateAccessToken(ctx echo.Context, app *app.Apps) (AuthResponse, error) {
	// Parse refresh token from request body
	var req LogoutRequest
	if err := ctx.Bind(&req); err != nil || req.RefreshToken == "" {
		return AuthResponse{}, utils.NewBadRequest("refresh token is required")
	}

	token, err := utils.ValidateToken(app, req.RefreshToken)
	if err != nil || !token.Valid {
		return AuthResponse{}, utils.NewForbidden("refresh token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AuthResponse{}, utils.NewBadRequest("invalid claims in refresh token")
	}

	data, ok := claims["data"].(map[string]interface{})
	if !ok {
		return AuthResponse{}, utils.NewBadRequest("Invalid data in claims")
	}

	userID, ok := data["id"].(string)
	if !ok {
		return AuthResponse{}, utils.NewBadRequest("invalid or missing user ID in token")
	}

	existingUser, err := a.repo.FindById(ctx, userID)
	if err != nil {
		return AuthResponse{}, utils.NewBadRequest("user not found")
	}

	// Aggregate permissions
	var allPermissions []string
	for _, role := range existingUser.RolesData {
		allPermissions = append(allPermissions, role.Permissions...)
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
		return AuthResponse{}, utils.NewInternal(err.Error())
	}

	return AuthResponse{
		Token:        newAccessToken,
		RefreshToken: req.RefreshToken,
		Data:         newPayload,
	}, nil
}
