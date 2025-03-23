package utils

import (
	"fmt"
	"time"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/golang-jwt/jwt/v5"
)

// var secretKey = []byte(configs.GeneralConfig.Security.JWT_SECRET_KEY)

// Function to create JWT tokens with claims
func CreateToken(config *config.Config, data interface{}) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(time.Hour * time.Duration(config.Security.JWTExpired)).Unix(),
		"iat":  time.Now().Unix(),
	})

	fmt.Printf("Token claims added: %+v\n", claims)

	tokenString, err := claims.SignedString("secretKey")
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Function to verify JWT tokens
func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return "secretKey", nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}
