package service

import (
	models "MITI_ART/Models"
	Utils "MITI_ART/Utils"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterVendor registers a new vendor using GORM transactions
func RegisterVendor(db *gorm.DB, VendorEmail string, VendorFirstName string, VendorOtherName string, VendorPassword string, role string, VendorTin int, ShopName string) (string, error) {
	// Convert role string to Role ENUM
	userRole := models.Role(strings.ToUpper(role))

	// Checking if user already exists
	var existingUser models.User
	if err := db.Where("email = ?", VendorEmail).First(&existingUser).Error; err == nil {
		return "", errors.New("user already exists")
	} else if err != gorm.ErrRecordNotFound {
		return "", errors.New("database error: " + err.Error())
	}

	// Hashing the password dukoresheje Utils.HashPassword
	hashedPassword, salt, err := Utils.HashPassword(VendorPassword)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// Begin transaction for multiple insertinons muri database yacu
	tx := db.Begin()
	if tx.Error != nil {
		return "", errors.New("failed to start transaction: " + tx.Error.Error())
	}

	// Rollback function in case of failure
	rollback := func() {
		tx.Rollback()
	}

	// Insert new user wumucuruzi
	newUser := models.User{
		FirstName: VendorFirstName,
		OtherName: VendorOtherName,
		Email:     VendorEmail,
		Password:  hashedPassword,
		Salt:      salt,
		Role:      userRole,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		rollback()
		return "", errors.New("failed to register user: " + err.Error())
	}

	// Insert vendor linked to user
	newVendor := models.Vendor{
		UserID:       newUser.ID,
		BusinessName: ShopName,
		TaxPin:       int64(VendorTin),
		Approved:     false,
	}

	if err := tx.Create(&newVendor).Error; err != nil {
		rollback()
		return "", errors.New("failed to register vendor: " + err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		rollback()
		return "", errors.New("failed to commit transaction: " + err.Error())
	}

	return "Vendor registered successfully", nil
}

// Registering Products for a vendor
func RegisterProduct(db *gorm.DB, VendorID uuid.UUID, ProductName string, ProductPrice float64, ProductCategory string, ProductMaterial string, ProductImageURL string) (string, error) {

	// Convert material string to Material ENUM
	productMaterial := models.Material(strings.ToUpper(ProductMaterial))

	// Insert product
	newProduct := models.Product{
		VendorID: VendorID,
		Name:     ProductName,
		Price:    ProductPrice,
		Category: ProductCategory,
		Material: productMaterial,
		ImageURL: ProductImageURL,
	}

	if err := db.Create(&newProduct).Error; err != nil {
		return "", errors.New("failed to register product: " + err.Error())
	}

	return "Product registered successfully", nil
}
