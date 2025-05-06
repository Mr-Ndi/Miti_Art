package utils

import (
	models "MITI_ART/Models"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetVendorIDByEmail(db *gorm.DB, email string) (uuid.UUID, error) {
	var vendor models.Vendor
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return uuid.Nil, errors.New("user not found")
	}
	if err := db.Where("user_id = ?", user.ID).First(&vendor).Error; err != nil {
		return uuid.Nil, errors.New("vendor not found for this user")
	}

	return vendor.UserID, nil
}

func UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	uploadDir := "uploads"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	fileType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(fileType, "image/") {
		return "", fmt.Errorf("only image files are allowed")
	}

	filePath := filepath.Join(uploadDir, fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename))

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return filePath, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserPhoneByID(db *gorm.DB, userID uuid.UUID) (string, error) {
	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return "", errors.New("user not found")
	}
	return user.Phone, nil
}
