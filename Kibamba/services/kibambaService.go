package services

import (
	models "MITI_ART/Models"
	utils "MITI_ART/Utils"
	"bytes"
	"fmt"
	"log"
	"os"

	"context"
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
func SeedAdmin(db *gorm.DB) {
	adminEmail := os.Getenv("ADMAIL")
	adminPassword := os.Getenv("ADPASSWORD")
	firstName := os.Getenv("FirstName")
	otherName := os.Getenv("OtherName")

	if adminEmail == "" || adminPassword == "" {
		log.Println("Missing credential environment variable. Admin user cannot be seeded.")
		return
	}

	var existingAdmin models.User
	err := db.Where("email = ?", adminEmail).First(&existingAdmin).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println("Database error while searching for admin:", err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		hashedPassword, salt, err := utils.HashPassword(adminPassword)
		if err != nil {
			panic("Failed to hash admin password: " + err.Error())
		}
		admin := models.User{
			FirstName: firstName,
			OtherName: otherName,
			Email:     adminEmail,
			Password:  hashedPassword,
			Salt:      salt,
			Role:      "ADMIN",
		}
		err = db.Create(&admin).Error
		if err != nil {
			panic("Failed to create admin user: " + err.Error())
		} else {
			log.Println("Admin user created successfully")
		}
	} else {
		log.Println("Admin user already exists")
	}
}

// Debugging helper function
func DebugHashPassword(password, salt string) string {
	saltBytes, _ := base64.StdEncoding.DecodeString(salt)
	hashed := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hashed)
}

// Login function
func Login(ctx context.Context, db *gorm.DB, email string, password string) (string, error) {
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			adminEmail := os.Getenv("ADMAIL")
			adminPassword := os.Getenv("ADPASSWORD")

			if email == adminEmail && password == adminPassword {
				SeedAdmin(db)
			}
			err := db.Where("email = ?", email).First(&user).Error
			if err != nil {
				return "", errors.New("failed, I mean creating admins")
			}
		} else {
			return "", errors.New("database error")
		}
	}
	if !checkPasswordHash(password, user.Password, user.Salt) {
		fmt.Println("Login failed: Wrong password")
		return "", errors.New("invalid credentials / Wrong password")
	}
	payload := map[string]interface{}{"email": user.Email, "role": user.Role, "user_id": user.ID.String()}
	token, err := utils.GenerateToken(payload)
	if err != nil {
		fmt.Println("Token generation failed:", err)
		return "", err
	}
	return token, nil
}

// Retriving all customers
func GetAllClients(db *gorm.DB) ([]models.User, error) {
	var clients []models.User
	if err := db.Where("role = ?", "CUSTOMER").Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

// Retriving all vendors
func GetAllVendors(db *gorm.DB) ([]models.Vendor, error) {
	var vendors []models.Vendor
	if err := db.Preload("User").Find(&vendors).Error; err != nil {
		return nil, err
	}
	return vendors, nil
}

// GetAllOrders - fetches all orders with related user and product data
func GetAllOrders(db *gorm.DB) ([]models.Order, error) {
	var orders []models.Order
	err := db.Preload("User").Preload("Product").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetAllProducts - fetches all products with vendor info
func GetAllProducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	err := db.Preload("Owner.User").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateVendor updates vendor details
func UpdateVendor(db *gorm.DB, id uuid.UUID, businessName string, taxPin int64) error {
	return db.Model(&models.Vendor{}).Where("user_id = ?", id).Updates(map[string]interface{}{
		"business_name": businessName,
		"tax_pin":       taxPin,
	}).Error
}

// UpdateClient updates client (user) details
func UpdateClient(db *gorm.DB, id uuid.UUID, fname, otherName, phone, city, address string) error {
	return db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"first_name": fname,
		"other_name": otherName,
		"phone":      phone,
		"city":       city,
		"address":    address,
	}).Error
}

// DeleteVendor removes a vendor and cascades to products if needed
func DeleteVendor(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&models.Vendor{}, "user_id = ?", id).Error
}

// DeleteClient removes a client (user) and cascades if needed
func DeleteClient(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&models.User{}, "id = ?", id).Error
}
