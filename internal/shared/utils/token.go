package utils

import (
	"fmt"
	"time"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/core/auth/model"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateAccessToken creates a JWT access token
func GenerateAccessToken(cfg *config.Config, payload interface{}) (string, error) {
	return createJWT(cfg.Security.JWTSecretKey, payload, time.Minute*time.Duration(cfg.Security.JWTExpired))
}

// GenerateRefreshToken creates a JWT refresh token
func GenerateRefreshToken(cfg *config.Config) (string, error) {
	return createJWT(cfg.Security.JWTSecretKey, nil, time.Hour*time.Duration(cfg.Security.JWTRefreshTokenExpired))
}

// createJWT generates a JWT token with a given expiration time
func createJWT(secretKey string, payload interface{}, expiration time.Duration) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": payload,
		"exp":  time.Now().Add(expiration).Unix(),
		"iat":  time.Now().Unix(),
	})

	return claims.SignedString([]byte(secretKey))
}

// ValidateToken verifies the given JWT token
func ValidateToken(cfg *config.Config, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Security.JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

// GenerateAuthToken creates both access and refresh tokens
func GenerateAuthToken(cfg *config.Config, payload interface{}) (model.AuthResponse, error) {
	accessToken, err := GenerateAccessToken(cfg, payload)
	if err != nil {
		return model.AuthResponse{}, err
	}

	refreshToken, err := GenerateRefreshToken(cfg)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		Data:         payload,
	}, nil
}
