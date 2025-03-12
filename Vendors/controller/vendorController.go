package controller

import (
	utils "MITI_ART/Utils"
	"MITI_ART/Vendors/service"
	"MITI_ART/prisma/miti_art"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHandle handles function
func RegisterHandle(c *gin.Context, prisma *miti_art.PrismaClient) {

	vendorToken := c.GetHeader("Authorization")
	if vendorToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorised Token is required"})
		return
	}

	// extrating the payload if token is valid
	payload, error := utils.ValidateToken(vendorToken)
	if error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid or expired token"})
		return
	}

	// gettin what we nees from the payload
	VendorEmail, emailOk := payload["vendorEmail"].(string)
	VendorFirstName, firstNameOk := payload["VendorFirstName"].(string)
	VendorOtherName, otheNameOk := payload["VendorOtherName"].(string)

	if !emailOk || !firstNameOk || !otheNameOk {
		fmt.Println("Vendor quirements available aren't enought")
	}
	// Request body
	var req struct {
		VendorPassword string `json:"vendorPassword" binding:"required"`
		VendorTin      int    `json:"vendorTin" binding:"required"`
	}

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Call service function
	message, err := service.RegisterVendor(prisma, VendorEmail, VendorFirstName, VendorOtherName, req.VendorPassword, req.VendorTin)

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
		"message":      message,
		"Vendor email": VendorEmail,
	})
}
