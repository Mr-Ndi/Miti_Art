package controller

import (
	"MITI_ART/src/Client/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =============================
// Controller Functions
// =============================

// RegisterHandle godoc
// @Summary Register a new client
// @Description Registers a new client with email, name, and password
// @Tags client
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Client registration input"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clients/register [post]
func RegisterHandle(c *gin.Context, db *gorm.DB) {
	var req struct {
		ClientFirstName string `json:"clientFirstName" binding:"required"`
		ClientOtherName string `json:"clientOtherName" binding:"required"`
		ClientEmail     string `json:"clientEmail" binding:"required,email"`
		ClientPassword  string `json:"clientPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	message, err := service.RegisterClient(db, req.ClientEmail, req.ClientFirstName, req.ClientOtherName, req.ClientPassword)
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": message,
		"email":   req.ClientEmail,
	})
}

// GetFurniture godoc
// @Summary Get all furniture
// @Description Fetches all available furniture products
// @Tags client
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /furniture [get]
func GetFurniture(c *gin.Context, db *gorm.DB) {
	products, err := service.Products(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

// GetFurnitureDetails godoc
// @Summary Get furniture by ID
// @Description Fetches details for a specific furniture item
// @Tags client
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /furniture/{id} [get]
func GetFurnitureDetails(c *gin.Context, db *gorm.DB) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	product, err := service.Product(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// CreateOrder godoc
// @Summary Create an order
// @Description Places a new order for a given product and quantity
// @Tags client
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Order request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /orders [post]
func CreateOrder(c *gin.Context, db *gorm.DB) {
	var req struct {
		ProductID uuid.UUID `json:"productID" binding:"required"`
		Quantity  int       `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDAny.(uuid.UUID)

	orderID, amount, message, err := service.HandleOrder(db, req.ProductID, req.Quantity, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": message,
		"orderID": orderID,
		"amount":  amount,
	})
}

// AppendWishList godoc
// @Summary Add to wishlist
// @Description Adds a product to the authenticated user's wishlist
// @Tags client
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Wishlist request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /wishlist [post]
func AppendWishList(c *gin.Context, db *gorm.DB) {
	var req struct {
		ProductID uuid.UUID `json:"productID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDAny.(uuid.UUID)

	id, message, err := service.WishList(db, req.ProductID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": message,
		"orderID": id,
	})
}

// ListUserOrders godoc
// @Summary List user orders
// @Description Gets a list of a user's orders with optional status and date filters
// @Tags client
// @Produce json
// @Param userID path string true "User ID"
// @Param status query string false "Order status"
// @Param from query string false "Start date (RFC3339 format)"
// @Param to query string false "End date (RFC3339 format)"
// @Success 200 {array} map[string]interface{}
// @Router /orders/user/{userID} [get]
func ListUserOrders(c *gin.Context, db *gorm.DB) {
	userID := c.Param("userID")
	status := c.Query("status")
	start := c.Query("from")
	end := c.Query("to")

	uuid := uuid.MustParse(userID)
	startTime, _ := time.Parse(time.RFC3339, start)
	endTime, _ := time.Parse(time.RFC3339, end)

	orders, err := service.GetUserOrders(db, uuid, &status, &startTime, &endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
