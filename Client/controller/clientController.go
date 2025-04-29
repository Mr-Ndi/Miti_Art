package controller

import (
	"MITI_ART/Client/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterHandle handles function
func RegisterHandle(c *gin.Context, db *gorm.DB) {
	// Request body
	var req struct {
		ClientFirstName string `json:"clientFirstName" binding:"required"`
		ClientOtherName string `json:"clientOtherName" binding:"required"`
		ClientEmail     string `json:"clientEmail" binding:"required,email"`
		ClientPassword  string `json:"clientPassword" binding:"required"`
	}

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Call service function
	message, err := service.RegisterClient(db, req.ClientEmail, req.ClientFirstName, req.ClientOtherName, req.ClientPassword)

	// Handle service response
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
		"email":   req.ClientEmail,
	})
}

// Using Furniture finder function in service
func GetFurniture(c *gin.Context, db *gorm.DB) {
	products, err := service.Products(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

// Using single Furniture finder function in service
func GetFurnitureDetails(c *gin.Context, db *gorm.DB) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}
	products, err := service.Product(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

// Using id to place order while logged in
func CreateOrder(c *gin.Context, db *gorm.DB) {
	var req struct {
		ProductID uuid.UUID `gorm:"not null"`
		Quantity  int       `gorm:"not null"`
		UserID    uuid.UUID `gorm:"not null"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request:" + err.Error()})
		return
	}
}
