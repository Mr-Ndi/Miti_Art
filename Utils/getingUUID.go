package utils

import (
	models "MITI_ART/Models"
	"errors"

	"gorm.io/gorm"
)

func GetVendorIDByEmail(db *gorm.DB, email string) (string, error) {
	var vendor models.Vendor

	// First, find the user with the given email
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	// Then, use the user ID to find the vendor
	if err := db.Where("user_id = ?", user.ID).First(&vendor).Error; err != nil {
		return "", errors.New("vendor not found for this user")
	}

	return vendor.UserID.String(), nil
}
