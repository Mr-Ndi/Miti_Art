package controller

import (
	"MITI_ART/Client/service"
	"MITI_ART/prisma/miti_art"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHandle handles function
func RegisterHandle(c *gin.Context, prisma *miti_art.PrismaClient) {
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
	message, err := service.RegisterClient(prisma, req.ClientEmail, req.ClientFirstName, req.ClientOtherName, req.ClientPassword)

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
