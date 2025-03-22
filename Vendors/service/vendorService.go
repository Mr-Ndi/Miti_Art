package service

import (
	Utils "MITI_ART/Utils"
	"MITI_ART/prisma/miti_art"
	"context"
	"errors"
	"fmt"
	"strings"
)

// RegisterVendor registers a new vendor with manual transaction handling
func RegisterVendor(prisma *miti_art.PrismaClient, VendorEmail string, VendorFirstName string, VendorOtherName string, VendorPassword string, role string, VendorTin int, ShopName string) (string, error) {
	ctx := context.Background()

	// Check if user already exists
	existingUser, err := prisma.User.FindUnique(
		miti_art.User.Email.Equals(VendorEmail),
	).Exec(ctx)

	if err != nil {
		return "", errors.New("database error: " + err.Error())
	}

	if existingUser == nil {
		fmt.Println("No existing user found with email:", VendorEmail) // Debugging log
	} else {
		fmt.Println("Existing user found:", existingUser.ID)
	}

	// Hash the password
	hashedPassword, salt, err := Utils.HashPassword(VendorPassword)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// Ensure the role is correctly set as an ENUM
	userRole := miti_art.Role(strings.ToUpper(role))

	// Start transaction
	_, err = prisma.Prisma.ExecuteRaw("BEGIN").Exec(ctx)
	if err != nil {
		return "", errors.New("failed to start transaction: " + err.Error())
	}

	// Rollback function in case of failure
	rollback := func() {
		prisma.Prisma.ExecuteRaw("ROLLBACK").Exec(ctx)
	}

	// Insert user
	newUser, err := prisma.User.CreateOne(
		miti_art.User.FirstName.Set(VendorFirstName),
		miti_art.User.OtherName.Set(VendorOtherName),
		miti_art.User.Email.Set(VendorEmail),
		miti_art.User.Password.Set(hashedPassword),
		miti_art.User.Salt.Set(salt),
		miti_art.User.Role.Set(userRole),
	).Exec(ctx)

	if err != nil {
		rollback()
		return "", errors.New("failed to register user: " + err.Error())
	}

	// Insert vendor linked to user
	_, err = prisma.Vendor.CreateOne(
		miti_art.Vendor.User.Link(miti_art.User.ID.Equals(newUser.ID)),
		miti_art.Vendor.BusinessName.Set(ShopName),
		miti_art.Vendor.TaxPin.Set(VendorTin),
		miti_art.Vendor.Approved.Set(false),
	).Exec(ctx)

	if err != nil {
		rollback()
		return "", errors.New("failed to register vendor: " + err.Error())
	}

	// Commit transaction
	_, err = prisma.Prisma.ExecuteRaw("COMMIT").Exec(ctx)
	if err != nil {
		rollback()
		return "", errors.New("failed to commit transaction: " + err.Error())
	}

	return "Vendor registered successfully", nil
}
