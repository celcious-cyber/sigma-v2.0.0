package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"sigma-api/internal/shared/config"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"` // "admin" or "student"
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT for a user or student
func GenerateToken(userID uint, role string, loginType string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		Type:   loginType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// ValidateToken parses and validates a JWT string
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
