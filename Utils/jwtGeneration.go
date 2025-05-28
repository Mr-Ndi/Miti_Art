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
		// fmt.Println(".env file loaded successfully.")
	}

	secret = os.Getenv("SECRET_KEY")

	if secret == "" {
		fmt.Println("SECRET_KEY is not set! Check your .env file.")
	} else {
		// fmt.Println("SECRET_KEY loaded successfully:", secret)
	}
}

// Generate JWT token
func GenerateToken(payload map[string]interface{}) (string, error) {
	if secret == "" {
		return "", errors.New("SECRET_KEY is not set")
	}

	expiration := time.Now().Add(time.Hour).Unix()
	payload["exp"] = expiration

	claims := jwt.MapClaims{}

	for key, value := range payload {
		claims[key] = value
	}

	claims["exp"] = time.Now().Add(time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// Validate the JWT token
func ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims structure")
	}

	return claims, nil

}
