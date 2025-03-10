package service

import (
	utils "MITI_ART/Utils"
	"MITI_ART/prisma/miti_art"
	"context"
	"errors"
)

// RegisterClient registers a new client (customer)
func RegisterClient(prisma *miti_art.PrismaClient, ClientEmail string, ClientFirstName string, ClientOtherName string, ClientPassword string) error {
	ctx := context.Background()

	// Check if the user already exists
	existingUser, err := prisma.User.FindUnique(
		miti_art.User.Email.Equals(ClientEmail),
	).Exec(ctx)

	if err != nil {
		return errors.New("database error: " + err.Error())
	}

	if existingUser != nil {
		return errors.New("email already registered")
	}

	// Hash the password user yatanze
	hashedPassword, salt, err := utils.HashPassword(ClientPassword)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	// Create new user as a "customer" role muri database
	_, err = prisma.User.CreateOne(
		miti_art.User.FirstName.Set(ClientFirstName),
		miti_art.User.OtherName.Set(ClientOtherName),
		miti_art.User.Email.Set(ClientEmail),
		miti_art.User.Password.Set(hashedPassword),
		miti_art.User.Salt.Set(salt),
		miti_art.User.Role.Set("customer"),
	).Exec(ctx)

	if err != nil {
		return errors.New("failed to register user: " + err.Error())
	}

	return nil
}
