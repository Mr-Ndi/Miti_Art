package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var secret string

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	} else {
		fmt.Println(".env file loaded successfully.")
	}

	secret = os.Getenv("SECRET_KEY")

	if secret == "" {
		fmt.Println("SECRET_KEY is not set! Check your .env file.")
	} else {
		fmt.Println("SECRET_KEY loaded successfully:", secret)
	}
}

// Generate JWT token
func GenerateToken(payload []string) (string, error) {
	if secret == "" {
		return "", errors.New("SECRET_KEY is not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString([]byte(secret))
}

// Validate the JWT token
func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if payload, exists := claims["payload"].(string); exists {
			return payload, nil
		}
		return "", errors.New("invalid token payload")
	}

	return "", errors.New("invalid token")
}
