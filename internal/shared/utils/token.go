package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HasanNugroho/starter-golang/internal/app"
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
func ValidateToken(app *app.Apps, tokenStr string) (*jwt.Token, error) {
	if app.Config.Redis.Enabled && IsTokenRevoked(app, tokenStr) {
		return nil, fmt.Errorf("token has been revoked")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Config.Security.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return token, nil
}

func GenerateAuthToken(app *app.Apps, payload interface{}) (accessToken string, refreshToken string, err error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	var parsedMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &parsedMap); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	accessToken, err = createJWT(app.Config.Security.JWTSecretKey, parsedMap, time.Minute*time.Duration(app.Config.Security.JWTExpired))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = createJWT(app.Config.Security.JWTSecretKey, nil, time.Hour*time.Duration(app.Config.Security.JWTRefreshTokenExpired))
	if err != nil {
		return "", "", err
	}

	userID, ok := parsedMap["id"].(string)
	if !ok || userID == "" {
		return "", "", fmt.Errorf("user ID not found or invalid")
	}

	// Store refresh token
	if err := StoreRefreshToken(app, userID, refreshToken, time.Hour*time.Duration(app.Config.Security.JWTRefreshTokenExpired)); err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// RefreshAccessToken validates refresh token and returns a new access token
func RefreshAccessToken(app *app.Apps, refreshToken string) (string, error) {
	ctx := context.Background()

	// Cek apakah token valid
	_, err := ValidateToken(app, refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Cek apakah refresh token sudah tidak berlaku
	key := "refresh_token:" + refreshToken
	userID, err := app.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("refresh token not found or revoked")
	}

	newPayload := map[string]interface{}{"user_id": userID}
	newAccessToken, err := createJWT(app.Config.Security.JWTSecretKey, newPayload, time.Minute*time.Duration(app.Config.Security.JWTExpired))
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

// RevokeToken stores the token in Redis with an expiration time
func RevokeToken(app *app.Apps, tokenString string, refreshToken string) error {
	redisClient := app.Redis
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
		if err := RevokeRefreshToken(app, refreshToken); err != nil {
			return fmt.Errorf("failed to revoke token")
		}
	}

	return nil
}

// RevokeRefreshToken deletes refresh token from Redis
func RevokeRefreshToken(app *app.Apps, refreshToken string) error {
	key := "refresh_token:" + refreshToken
	return app.Redis.Del(context.Background(), key).Err()
}

// IsTokenRevoked checks if a token is in the Redis blacklist
func IsTokenRevoked(app *app.Apps, tokenString string) bool {
	redisClient := app.Redis
	_, err := redisClient.Get(context.Background(), "blacklist:"+tokenString).Result()
	return err == nil
}

func StoreRefreshToken(app *app.Apps, userID string, refreshToken string, expiration time.Duration) error {
	key := "refresh_token:" + refreshToken
	return app.Redis.Set(context.Background(), key, userID, expiration).Err()
}
