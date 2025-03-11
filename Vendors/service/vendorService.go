package service

import (
	Utils "MITI_ART/Utils"
	"MITI_ART/prisma/miti_art"
	"context"
	"errors"
)

// RegisterClient registers a new client
func RegisterClient(prisma *miti_art.PrismaClient, ClientEmail string, ClientFirstName string, ClientOtherName string, ClientPassword string) (string, error) {
	ctx := context.Background()

	// Check if user already exists
	existingUser, err := prisma.User.FindUnique(
		miti_art.User.Email.Equals(ClientEmail),
	).Exec(ctx)

	if err != nil && err.Error() != "ErrNotFound" {
		return "", errors.New("database error: " + err.Error())
	}

	if existingUser != nil {
		return "", errors.New("user with that email already registered")
	}

	// Hash the password
	hashedPassword, salt, err := Utils.HashPassword(ClientPassword)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// Create new user
	_, err = prisma.User.CreateOne(
		miti_art.User.FirstName.Set(ClientFirstName),
		miti_art.User.OtherName.Set(ClientOtherName),
		miti_art.User.Email.Set(ClientEmail),
		miti_art.User.Password.Set(hashedPassword),
		miti_art.User.Salt.Set(salt),
		miti_art.User.Role.Set("customer"),
	).Exec(ctx)

	if err != nil {
		return "", errors.New("failed to register user: " + err.Error())
	}

	return "User registered successfully", nil
}
