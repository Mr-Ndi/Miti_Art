package service

import (
	Utils "MITI_ART/Utils"
	"MITI_ART/prisma/miti_art"
	"context"
	"errors"
)

// RegisterClient registers a new client
func RegisterVendor(prisma *miti_art.PrismaClient, VendorEmail string, VendorFirstName string, VendorOtherName string, VendorPassword string, role string, VendorTin int) (string, error) {
	ctx := context.Background()

	// Check if user already exists
	existingUser, err := prisma.User.FindUnique(
		miti_art.User.Email.Equals(VendorEmail),
	).Exec(ctx)

	if err != nil && err.Error() != "ErrNotFound" {
		return "", errors.New("database error: " + err.Error())
	}

	if existingUser != nil {
		return "", errors.New("user with that email already registered")
	}

	// Hash the password
	hashedPassword, salt, err := Utils.HashPassword(VendorPassword)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// Create new user
	_, err = prisma.User.CreateOne(
		miti_art.User.FirstName.Set(VendorFirstName),
		miti_art.User.OtherName.Set(VendorOtherName),
		miti_art.User.Email.Set(VendorEmail),
		miti_art.User.Password.Set(hashedPassword),
		miti_art.User.Salt.Set(salt),
		miti_art.User.Role.Set(role),
	).Exec(ctx)

	if err != nil {
		return "", errors.New("failed to register user: " + err.Error())
	}

	return "Vendor registered successfully", nil
}
