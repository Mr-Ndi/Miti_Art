package controller

import (
	"MITI_ART/src/Client/dto"
	"MITI_ART/src/Client/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterHandle godoc
// @Summary Register a new client
// @Description Registers a new client with email, name, and password
// @Tags client
// @Accept json
// @Produce json
// @Param body body dto.ClientRegisterRequest true "Client registration input"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router //user/register [post]
func RegisterHandle(c *gin.Context, db *gorm.DB) {
	var req dto.ClientRegisterRequest
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
// @Router /user/furniture [get]
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
// @Router /user/furniture/{id} [get]
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
// @Description Places a new order for a given product and quantity (authentication required)
// @Tags client
// @Accept json
// @Produce json
// @Param body body dto.CreateOrderRequest true "Order request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /user/orders [post]
func CreateOrder(c *gin.Context, db *gorm.DB) {
	var req dto.OrderRequest
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
// @Description Adds a product to the authenticated user's wishlist (authentication required)
// @Tags client
// @Accept json
// @Produce json
// @Param body body dto.WishListRequest true "Wishlist request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /user/wishlist [post]
func AppendWishList(c *gin.Context, db *gorm.DB) {
	var req dto.WishlistRequest
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
// @Description Gets a list of a user's orders with optional status and date filters (authentication required)
// @Tags client
// @Produce json
// @Param userID path string true "User ID"
// @Param status query string false "Order status"
// @Param from query string false "Start date (RFC3339 format)"
// @Param to query string false "End date (RFC3339 format)"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/orders/user/{userID} [get]
func ListUserOrders(c *gin.Context, db *gorm.DB) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID := c.Param("userID")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	status := c.Query("status")
	var startTime, endTime *time.Time

	if start := c.Query("from"); start != "" {
		t, err := time.Parse(time.RFC3339, start)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' date"})
			return
		}
		startTime = &t
	}

	if end := c.Query("to"); end != "" {
		t, err := time.Parse(time.RFC3339, end)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'to' date"})
			return
		}
		endTime = &t
	}

	orders, err := service.GetUserOrders(db, userUUID, &status, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
