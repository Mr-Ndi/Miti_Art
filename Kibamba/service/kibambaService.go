package services

import (
	"MITI_ART/prisma/miti_art"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

var jwtSecret = []byte("your_secret_key") // Replace with a secure secret

// Generate a random salt (16 bytes)
func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// Hash password using Argon2
func hashPassword(password string) (string, string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", "", err
	}

	saltBytes, _ := base64.StdEncoding.DecodeString(salt)
	hashed := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	hash := base64.StdEncoding.EncodeToString(hashed)

	return hash, salt, nil
}

// Verify password using Argon2
func checkPasswordHash(password, hash, salt string) bool {
	saltBytes, _ := base64.StdEncoding.DecodeString(salt)
	hashedAttempt := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	expectedHash, _ := base64.StdEncoding.DecodeString(hash)

	return string(hashedAttempt) == string(expectedHash)
}

// Generate JWT token
func generateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})
	return token.SignedString(jwtSecret)
}

// Seed an admin user (Creates if absent)
func SeedAdmin(prisma *miti_art.PrismaClient) {
	ctx := context.Background()
	adminEmail := "admin@example.com"   // Replace with actual admin email
	adminPassword := "AdminPassword123" // Replace with a secure password

	existingAdmin, _ := prisma.User.FindUnique(
		miti_art.User.Email.Equals(adminEmail),
	).Exec(ctx)

	if existingAdmin == nil {
		hashedPassword, salt, _ := hashPassword(adminPassword)
		prisma.User.CreateOne(
			miti_art.User.Email.Set(adminEmail),
			miti_art.User.Password.Set(hashedPassword),
			miti_art.User.Salt.Set(salt),
			miti_art.User.Role.Set("admin"),
		).Exec(ctx)
	}
}

// Login function
func Login(ctx context.Context, prisma *miti_art.PrismaClient, email, password string) (string, error) {
	user, err := prisma.User.FindUnique(
		miti_art.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if !checkPasswordHash(password, user.Password, user.Salt) {
		return "", errors.New("invalid credentials/Wrong password")
	}

	return generateToken(user.Email)
}
