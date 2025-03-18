package service

import (
	Utils "MITI_ART/Utils"
	"MITI_ART/prisma/miti_art"
	"context"
	"errors"
	"strings"
)

// RegisterVendor a function for registers a new vendor
func RegisterVendor(prisma *miti_art.PrismaClient, VendorEmail string, VendorFirstName string, VendorOtherName string, VendorPassword string, role string, VendorTin int, ShopName string) (string, error) {
	ctx := context.Background()

	// Check if user already exists
	existingUser, err := prisma.User.FindUnique(
		miti_art.User.Email.Equals(VendorEmail),
	).Exec(ctx)

	if err != nil {
		return "", errors.New("database error: " + err.Error())
	}

	if existingUser != nil {
		return "", errors.New("user with that email already registered")
	}

	// Hashing the password
	hashedPassword, salt, err := Utils.HashPassword(VendorPassword)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// Ensure the role is correctly set as an ENUM
	userRole := miti_art.Role(strings.ToUpper(role))

	// Create new user
	newUser, err := prisma.User.CreateOne(
		miti_art.User.FirstName.Set(VendorFirstName),
		miti_art.User.OtherName.Set(VendorOtherName),
		miti_art.User.Email.Set(VendorEmail),
		miti_art.User.Password.Set(hashedPassword),
		miti_art.User.Salt.Set(salt),
		miti_art.User.Role.Set(userRole),
	).Exec(ctx)

	if err != nil {
		return "", errors.New("failed to register user: " + err.Error())
	} else if newUser.ID == "" {
		return "", errors.New("failed to retrieve user ID for vendor creation")

	}

	// Creating a user record inside Vendor table
	_, err = prisma.Vendor.CreateOne(
		miti_art.Vendor.User.Link(miti_art.User.ID.Equals(newUser.ID)),
		miti_art.Vendor.BusinessName.Set(ShopName),
		miti_art.Vendor.TaxPin.Set(VendorTin),
		miti_art.Vendor.Approved.Set(false),
	).Exec(ctx)

	if err != nil {
		return "", errors.New("failed to register vendor: " + err.Error())
	}

	return "Vendor registered successfully", nil
}
