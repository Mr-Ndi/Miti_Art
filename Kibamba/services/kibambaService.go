package services

import (
	utils "MITI_ART/Utils"
	"bytes"
	"fmt"
	"os"

	"MITI_ART/prisma/miti_art"
	"context"
	"encoding/base64"
	"errors"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/argon2"
)

func init() {
	_ = godotenv.Load()
}

// Verify password using Argon2
func checkPasswordHash(password, hash, salt string) bool {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		fmt.Println("Error decoding salt:", err)
		return false
	}

	hashedAttempt := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	expectedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		fmt.Println("Error decoding stored hash:", err)
		return false
	}

	match := bytes.Equal(hashedAttempt, expectedHash)
	return match
}

// Seed an admin user (Creates if absent)
func SeedAdmin(prisma *miti_art.PrismaClient) {
	ctx := context.Background()
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	firstName := os.Getenv("FirstName")
	otherName := os.Getenv("OtherName")

	existingAdmin, err := prisma.User.FindUnique(
		miti_art.User.Email.Equals(adminEmail),
	).Exec(ctx)

	if err != nil || existingAdmin == nil {
		hashedPassword, salt, err := utils.HashPassword(adminPassword)
		if err != nil {
			panic("Failed to hash admin password: " + err.Error())
		}

		_, err = prisma.User.CreateOne(
			miti_art.User.FirstName.Set(firstName),
			miti_art.User.OtherName.Set(otherName),
			miti_art.User.Email.Set(adminEmail),
			miti_art.User.Password.Set(hashedPassword),
			miti_art.User.Salt.Set(salt),
			miti_art.User.Role.Set("admin"),
		).Exec(ctx)

		if err != nil {
			panic("Failed to create admin user: " + err.Error())
		}
	}
}

// Debugging helper function
func DebugHashPassword(password, salt string) string {
	saltBytes, _ := base64.StdEncoding.DecodeString(salt)
	hashed := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hashed)
}

// Login function
func Login(ctx context.Context, prisma *miti_art.PrismaClient, email string, password string) (string, error) {
	user, err := prisma.User.FindUnique(
		miti_art.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil || user == nil {
		adminEmail := os.Getenv("ADMIN_EMAIL")
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if email == adminEmail && password == adminPassword {
			SeedAdmin(prisma)
		}

		user, err = prisma.User.FindUnique(
			miti_art.User.Email.Equals(email),
		).Exec(ctx)

		if err != nil || user == nil {
			return "", errors.New("failed, contact admins")
		}
	}

	if !checkPasswordHash(password, user.Password, user.Salt) {
		fmt.Println("Login failed: Wrong password")
		return "", errors.New("invalid credentials / Wrong password")
	}

	payload := map[string]interface{}{"email": user.Email, "role": user.Role}
	token, err := utils.GenerateToken(payload)
	if err != nil {
		fmt.Println("Token generation failed:", err)
		return "", err
	}

	fmt.Println("✅ Token generated successfully:", token)

	return token, nil
}
