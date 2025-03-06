package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Load environment variables
func init() {
	_ = godotenv.Load() // Load .env file if available
}

// Secret key for JWT
var secret = os.Getenv("SECRET_KEY")

// Generate JWT token
func GenerateToken(payload []string) (string, error) {
	if secret == "" {
		return "", errors.New("SECRET_KEY is not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Hour).Unix(), // Token expires in 1 hour
	})

	return token.SignedString([]byte(secret))
}

// Validate the JWT token
func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Ensure signing method is correct
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	// Extract payload from claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if payload, exists := claims["payload"].(string); exists {
			return payload, nil
		}
		return "", errors.New("invalid token payload")
	}

	return "", errors.New("invalid token")
}
