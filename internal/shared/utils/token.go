package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/golang-jwt/jwt/v5"
)

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
func ValidateToken(config *config.Config, tokenStr string) (*jwt.Token, error) {
	if config.Redis.Enabled && IsTokenRevoked(config, tokenStr) {
		return nil, fmt.Errorf("token has been revoked")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Security.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return token, nil
}

func GenerateAuthToken(config *config.Config, payload interface{}) (accessToken string, refreshToken string, err error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	var parsedMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &parsedMap); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	accessToken, err = createJWT(config.Security.JWTSecretKey, parsedMap, time.Minute*time.Duration(config.Security.JWTExpired))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = createJWT(config.Security.JWTSecretKey, nil, time.Hour*time.Duration(config.Security.JWTRefreshTokenExpired))
	if err != nil {
		return "", "", err
	}

	userID, ok := parsedMap["id"].(string)
	if !ok || userID == "" {
		return "", "", fmt.Errorf("user ID not found or invalid")
	}

	// Store refresh token
	if err := StoreRefreshToken(config, userID, refreshToken, time.Hour*time.Duration(config.Security.JWTRefreshTokenExpired)); err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// RefreshAccessToken validates refresh token and returns a new access token
func RefreshAccessToken(config *config.Config, refreshToken string) (string, error) {
	ctx := context.Background()

	// Cek apakah token valid
	_, err := ValidateToken(config, refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Cek apakah refresh token sudah tidak berlaku
	key := "refresh_token:" + refreshToken
	userID, err := config.Redis.Client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("refresh token not found or revoked")
	}

	newPayload := map[string]interface{}{"user_id": userID}
	newAccessToken, err := createJWT(config.Security.JWTSecretKey, newPayload, time.Minute*time.Duration(config.Security.JWTExpired))
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

// RevokeToken stores the token in Redis with an expiration time
func RevokeToken(config *config.Config, tokenString string, refreshToken string) error {
	redisClient := config.Redis.Client
	ctx := context.Background()

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("invalid expiration claim")
	}

	ttl := time.Until(time.Unix(int64(exp), 0))
	if err = redisClient.Set(ctx, "blacklist:"+tokenString, "revoked", ttl).Err(); err != nil {
		return err
	}

	if refreshToken != "" {
		if err := RevokeRefreshToken(config, refreshToken); err != nil {
			return fmt.Errorf("failed to revoke token")
		}
	}

	return nil
}

// RevokeRefreshToken deletes refresh token from Redis
func RevokeRefreshToken(config *config.Config, refreshToken string) error {
	key := "refresh_token:" + refreshToken
	return config.Redis.Client.Del(context.Background(), key).Err()
}

// IsTokenRevoked checks if a token is in the Redis blacklist
func IsTokenRevoked(config *config.Config, tokenString string) bool {
	redisClient := config.Redis.Client
	_, err := redisClient.Get(context.Background(), "blacklist:"+tokenString).Result()
	return err == nil
}

func StoreRefreshToken(config *config.Config, userID string, refreshToken string, expiration time.Duration) error {
	key := "refresh_token:" + refreshToken
	return config.Redis.Client.Set(context.Background(), key, userID, expiration).Err()
}
