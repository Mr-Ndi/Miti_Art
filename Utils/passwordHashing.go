package utils

import (
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

// Hash password using Argon2
func HashPassword(password string) (string, string, error) {
	saltBytes, err := GenerateSalt()
	if err != nil {
		return "", "", err
	}

	hashed := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	hash := base64.StdEncoding.EncodeToString(hashed)
	salt := base64.StdEncoding.EncodeToString(saltBytes)

	return hash, salt, nil
}
