package controllers

import (
	"context"
	"log"
	"net/http"
	"os"

	"MITI_ART/Kibamba/services"
	utils "MITI_ART/Utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func LoginHandler(c *gin.Context, db *gorm.DB) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := services.Login(context.Background(), db, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

type InviteRequest struct {
	VendorEmail     string `json:"VendorEmail" binding:"required"`
	VendorFirstName string `json:"VendorFirstName" binding:"required"`
	VendorOtherName string `json:"VendorOtherName" binding:"required"`
}

func InvitationHandler(c *gin.Context) {
	userEmail, exist := c.Get("userEmail")
	_ = userEmail //Guhagarika Complains za interpreter
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Acces"})
		return
	}
	var req InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data!"})
		return
	}
	payload := map[string]interface{}{"VendorEmail": req.VendorEmail, "VendorFirstName": req.VendorFirstName, "VendorOtherName": req.VendorOtherName, "role": "Vendor"}
	token, err := utils.GenerateToken(payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	result := utils.Invite(req.VendorEmail, req.VendorFirstName, req.VendorOtherName, token)

	c.JSON(http.StatusOK, gin.H{
		"status":  result,
		"message": "Invitation sent succesfully",
		"sent to": req.VendorEmail,
	})
}

func ViewClients(c *gin.Context, db *gorm.DB) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if userEmail != adminEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	clients, err := services.GetAllClients(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clients": clients})
}

func ViewVendors(c *gin.Context, db *gorm.DB) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if userEmail != adminEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	clients, err := services.GetAllClients(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vendors": clients})
}
