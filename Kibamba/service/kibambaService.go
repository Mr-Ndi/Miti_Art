package services

import (
	"context"
	"errors"
	"time"

	"your_project/prisma/db"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key")

// Hash password
func hashPassword(password string) (string, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Verify password
func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Generate JWT token
func generateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// Seed an admin user
func SeedAdmin(prisma *db.PrismaClient) {
	ctx := context.Background()
	adminEmail := "admin@example.com"
	adminPassword := "Admin@123"

	existingAdmin, _ := prisma.User.FindUnique(
		db.User.Email.Equals(adminEmail),
	).Exec(ctx)

	if existingAdmin == nil {
		hashedPassword, _ := hashPassword(adminPassword)
		prisma.User.CreateOne(
			db.User.Email.Set(adminEmail),
			db.User.Password.Set(string(hashedPassword)),
			db.User.Role.Set("admin"),
		).Exec(ctx)
	}
}

// Login function
func Login(ctx context.Context, prisma *db.PrismaClient, email, password string) (string, error) {
	user, _ := prisma.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if user == nil || !checkPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return generateToken(user.Email)
}
