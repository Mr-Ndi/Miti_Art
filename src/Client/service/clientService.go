package service

import (
	Utils "MITI_ART/Utils"
	utils "MITI_ART/Utils"
	models "MITI_ART/src/Models"
	"errors"
	"fmt"
	"time"

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

// Returning a product to the clients
func Product(db *gorm.DB, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := db.First(&product, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("product not found with ID: %s", id)
	}

	return &product, err
}

// Registering order
func HandleOrder(db *gorm.DB, productID uuid.UUID, quantity int, userID uuid.UUID) (uuid.UUID, float64, string, error) {
	phone, err := utils.GetUserPhoneByID(db, userID)
	if err != nil {
		return uuid.Nil, 0, "", err
	}

	var product models.Product
	if err := db.First(&product, "id = ?", productID).Error; err != nil {
		return uuid.Nil, 0, "", errors.New("product not found")
	}
	amount := product.Price * float64(quantity)

	newOrder := models.Order{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
	if err := db.Create(&newOrder).Error; err != nil {
		return uuid.Nil, 0, "", errors.New("failed to register order: " + err.Error())
	}

	token, err := utils.GetPaypackToken()
	if err != nil {
		return uuid.Nil, 0, "", errors.New("failed to get payment token")
	}

	err = utils.InitiatePayment(token, utils.PaymentPayload{
		Amount: amount,
		Number: phone,
	})
	if err != nil {
		return uuid.Nil, 0, "", errors.New("payment failed: " + err.Error())
	}

	return newOrder.ID, amount, "Order has been placed", nil
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

// Single endpoint with optional filters for Returning all products orders to the clients who ordered
func GetUserOrders(db *gorm.DB, userID uuid.UUID, status *string, startTime *time.Time, endTime *time.Time) ([]models.Order, error) {

	query := db.
		Preload("Product").
		Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if startTime != nil && endTime != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *startTime, *endTime)
	} else if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	} else if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	var orders []models.Order
	err := query.Find(&orders).Error
	return orders, err
}
