package controllers

import (
	"context"
	"net/http"

	"MITI_ART/Kibamba/services"
	utils "MITI_ART/Utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
