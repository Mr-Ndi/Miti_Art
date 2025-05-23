package controllers

import (
	"context"
	"log"
	"net/http"
	"os"

	"MITI_ART/Kibamba/services"
	utils "MITI_ART/Utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	vendors, err := services.GetAllVendors(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vendors"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"vendors": vendors})
}

func ViewOrders(c *gin.Context, db *gorm.DB) {
	userEmail, exists := c.Get("user_email")
	if !exists || userEmail != os.Getenv("ADMIN_EMAIL") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	orders, err := services.GetAllOrders(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func ViewAllProducts(c *gin.Context, db *gorm.DB) {
	userEmail, exists := c.Get("user_email")
	if !exists || userEmail != os.Getenv("ADMIN_EMAIL") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	products, err := services.GetAllProducts(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func EditVendor(c *gin.Context, db *gorm.DB) {
	var input struct {
		UserID       string `json:"user_id"`
		BusinessName string `json:"business_name"`
		TaxPin       int64  `json:"tax_pin"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
		return
	}

	err = services.UpdateVendor(db, id, input.BusinessName, input.TaxPin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vendor updated successfully"})
}

func EditClient(c *gin.Context, db *gorm.DB) {
	var input struct {
		UserID    string `json:"user_id"`
		FirstName string `json:"first_name"`
		OtherName string `json:"other_name"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
		City      string `json:"city"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	err = services.UpdateClient(db, id, input.FirstName, input.OtherName, input.Phone, input.City, input.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client updated successfully"})
}

func EliminateVendor(c *gin.Context, db *gorm.DB) {
	var input struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor ID"})
		return
	}

	err = services.DeleteVendor(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vendor eliminated successfully"})
}

func EliminateClient(c *gin.Context, db *gorm.DB) {
	var input struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	err = services.DeleteClient(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client eliminated successfully"})
}
