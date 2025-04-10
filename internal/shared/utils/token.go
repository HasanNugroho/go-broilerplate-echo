package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/golang-jwt/jwt/v5"
)

// createJWT generates a JWT token with a given expiration time
func createJWT(secretKey string, payload map[string]interface{}, expiration time.Duration) (string, error) {
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
		return nil, NewUnauthorized("Token is invalid or has been revoked")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Config.Security.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, NewUnauthorized("Token is invalid or has been revoked")
	}

	return token, nil
}

func GenerateAuthToken(app *app.Apps, payload interface{}) (accessToken string, refreshToken string, err error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", NewInternal("failed to generate token")
	}

	var parsedMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &parsedMap); err != nil {
		return "", "", NewInternal("failed to generate token")
	}

	accessToken, err = createJWT(app.Config.Security.JWTSecretKey, parsedMap, time.Hour*time.Duration(app.Config.Security.JWTExpired))
	if err != nil {
		return "", "", NewInternal("failed to generate token")
	}

	refreshToken, err = createJWT(app.Config.Security.JWTSecretKey, map[string]interface{}{"id": parsedMap["id"]}, time.Hour*time.Duration(app.Config.Security.JWTRefreshTokenExpired))
	if err != nil {
		return "", "", NewInternal("failed to generate token")
	}

	userID, ok := parsedMap["id"].(string)
	if !ok || userID == "" {
		return "", "", NewBadRequest("user ID not found or invalid")
	}

	// Store refresh token
	if err := StoreRefreshToken(app, userID, refreshToken, time.Hour*time.Duration(app.Config.Security.JWTRefreshTokenExpired)); err != nil {
		return "", "", NewInternal("failed to store refresh token")
	}

	return accessToken, refreshToken, nil
}

// RefreshAccessToken validates refresh token and returns a new access token
func RefreshAccessToken(app *app.Apps, refreshToken string, newPayload map[string]interface{}) (string, error) {
	ctx := context.Background()

	// Cek apakah token valid
	_, err := ValidateToken(app, refreshToken)
	if err != nil {
		return "", NewBadRequest("invalid refresh token")
	}

	// Cek apakah refresh token sudah tidak berlaku
	key := "refresh_token:" + refreshToken
	_, err = app.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", NewUnauthorized("refresh token not found or revoked")
	}

	newAccessToken, err := createJWT(app.Config.Security.JWTSecretKey, newPayload, time.Minute*time.Duration(app.Config.Security.JWTExpired))
	if err != nil {
		return "", NewInternal("failed to generate token")
	}

	return newAccessToken, nil
}

// RevokeToken stores the token in Redis with an expiration time
func RevokeToken(app *app.Apps, tokenString string, refreshToken string) error {
	redisClient := app.Redis
	ctx := context.Background()

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return NewInternal("failed to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return NewBadRequest("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return NewBadRequest("invalid expiration claim")
	}

	ttl := time.Until(time.Unix(int64(exp), 0))
	if err = redisClient.Set(ctx, "blacklist:"+tokenString, "revoked", ttl).Err(); err != nil {
		return NewInternal("failed to store token in blacklist")
	}

	if refreshToken != "" {
		if err := RevokeRefreshToken(app, refreshToken); err != nil {
			return NewInternal("failed to revoke token")
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
