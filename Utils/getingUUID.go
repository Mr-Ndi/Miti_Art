package utils

import (
	models "MITI_ART/Models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetVendorIDByEmail(db *gorm.DB, email string) (uuid.UUID, error) {
	var vendor models.Vendor

	// First, find the user with the given email
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return uuid.Nil, errors.New("user not found")
	}

	// Then, use the user ID to find the vendor
	if err := db.Where("user_id = ?", user.ID).First(&vendor).Error; err != nil {
		return uuid.Nil, errors.New("vendor not found for this user")
	}

	return vendor.UserID, nil
}
