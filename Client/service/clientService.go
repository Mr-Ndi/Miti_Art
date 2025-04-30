package service

import (
	models "MITI_ART/Models"
	Utils "MITI_ART/Utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterClient registers a new client
func RegisterClient(db *gorm.DB, ClientEmail string, ClientFirstName string, ClientOtherName string, ClientPassword string) (string, error) {

	// Check if user already exists
	var existingUser models.User

	if err := db.Where("email = ?", ClientEmail).First(&existingUser).Error; err == nil {
		return "", errors.New("user with that email already registered")
	} else if err != gorm.ErrRecordNotFound {
		return "", errors.New("database error: " + err.Error())
	}

	// Hash the password
	hashedPassword, salt, err := Utils.HashPassword(ClientPassword)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// Create new user
	newUser := models.User{
		FirstName: ClientFirstName,
		OtherName: ClientOtherName,
		Email:     ClientEmail,
		Password:  hashedPassword,
		Salt:      salt,
		Role:      "customer",
	}
	// inserting the given data in the database
	if err := db.Create(&newUser).Error; err != nil {
		return "", errors.New("failed to register user: " + err.Error())
	}

	return "User registered successfully", nil
}

// Returning all products to the clients
func Products(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	results := db.Find(&products)
	return products, results.Error
}

// Returning single products to the clients
func Product(db *gorm.DB, id uuid.UUID) ([]models.Product, error) {
	var product []models.Product
	results := db.Find(&product, "id = ?", id)
	return product, results.Error
}

// Registering order
func Order(db *gorm.DB, ProductID uuid.UUID, Quantity int, UserID uuid.UUID) (uuid.UUID, string, error) {
	// Inserting the new product in the database
	newOrder := models.Order{
		UserID:    UserID,
		ProductID: ProductID,
		Quantity:  Quantity,
	}

	if err := db.Create(&newOrder).Error; err != nil {
		return uuid.Nil, "", errors.New("failed to register product: " + err.Error())
	}
	return newOrder.ID, "Order has been Placed", nil
}

// Adding product on the wish list
func WishList(db *gorm.DB, ProductID uuid.UUID, UserID uuid.UUID) (uuid.UUID, string, error) {
	// Inserting the new product in the database to wishlist table
	newElement := models.Order{
		UserID:    UserID,
		ProductID: ProductID,
	}

	if err := db.Create(&newElement).Error; err != nil {
		return uuid.Nil, "", errors.New("failed to register product: " + err.Error())
	}
	return newElement.ID, "Order has been Placed", nil
}
