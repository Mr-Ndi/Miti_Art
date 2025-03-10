package controllers

import (
	"context"
	"net/http"

	"MITI_ART/Kibamba/services"
	utils "MITI_ART/Utils"
	"MITI_ART/prisma/miti_art"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context, prisma *miti_art.PrismaClient) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := services.Login(context.Background(), prisma, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

type InviteRequest struct {
	VentureEmail string `json:"ventureEmail" binding:"required"`
	VentureName  string `json:"ventureName" binding:"required"`
}

func InvitationHandler(c *gin.Context) {
	userEmail, exist := c.Get("userEmail")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Acces"})
		return
	}
	var req InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data!"})
		return
	}
	result := utils.Invite(req.VentureEmail, req.VentureName)

	c.JSON(http.StatusOK, gin.H{
		"status":  result,
		"message": "Invitation sent succesfully",
		"sent_By": userEmail,
	})
}
