package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate a random salt (16 bytes)
func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}
