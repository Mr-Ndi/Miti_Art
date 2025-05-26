package controllers

import (
	"context"
	"log"
	"net/http"
	"os"

	utils "MITI_ART/Utils"
	"MITI_ART/src/Kibamba/dto"
	"MITI_ART/src/Kibamba/services"

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

// LoginHandler godoc
// @Summary Login user
// @Description Authenticates user and returns a token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body map[string]string true "Login input"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /user/login [post]
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

// InvitationHandler godoc
// @Summary Send vendor invitation
// @Description Sends invitation token to a vendor
// @Tags admin
// @Accept json
// @Produce json
// @Param body body dto.InvitationInput true "Invitation input"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /invite [post]
func InvitationHandler(c *gin.Context) {
	var req dto.InvitationInput

	userEmail, exist := c.Get("userEmail")
	_ = userEmail
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Access"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data!"})
		return
	}

	payload := map[string]interface{}{
		"VendorEmail":     req.VendorEmail,
		"VendorFirstName": req.VendorFirstName,
		"VendorOtherName": req.VendorOtherName,
		"role":            "Vendor",
	}
	token, err := utils.GenerateToken(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	result := utils.Invite(req.VendorEmail, req.VendorFirstName, req.VendorOtherName, token)

	c.JSON(http.StatusOK, gin.H{
		"status":  result,
		"message": "Invitation sent successfully",
		"sent_to": req.VendorEmail,
	})
}

// ViewClients godoc
// @Summary View all clients
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Security BearerAuth
// @Router /admin/view-clients [get]
func ViewClients(c *gin.Context, db *gorm.DB) {
	userEmail, exists := c.Get("user_email")
	if !exists || userEmail != os.Getenv("ADMIN_EMAIL") {
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

// ViewVendors godoc
// @Summary View all vendors
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/view-vendors [get]
func ViewVendors(c *gin.Context, db *gorm.DB) {
	vendors, err := services.GetAllVendors(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vendors"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"vendors": vendors})
}

// ViewOrders godoc
// @Summary View all orders
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Security BearerAuth
// @Router /admin/view-orders [get]
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

// ViewAllProducts godoc
// @Summary View all products
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Security BearerAuth
// @Router /admin/view-products [get]
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

// EditVendor godoc
// @Summary Edit vendor info
// @Tags admin
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Vendor update input"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /admin/edit-vendor [post]
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

// EditClient godoc
// @Summary Edit client info
// @Tags admin
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Client update input"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /admin/edit-client [post]
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

// EliminateVendor godoc
// @Summary Delete a vendor
// @Tags admin
// @Accept json
// @Produce json
// @Param body body map[string]string true "Vendor ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /admin/eliminate-vendor [post]
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

// EliminateClient godoc
// @Summary Delete a client
// @Tags admin
// @Accept json
// @Produce json
// @Param body body map[string]string true "Client ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /admin/eliminate-client [post]
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
