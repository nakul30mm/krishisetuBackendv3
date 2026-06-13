package utils

import (
	"time"

	"krishisetu-backend/config"

	"github.com/golang-jwt/jwt/v5"
)

// generates access token for the user logging in
func GenerateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret())
}
