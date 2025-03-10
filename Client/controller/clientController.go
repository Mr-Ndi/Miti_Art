package controller

import (
	"MITI_ART/Client/service"
	"MITI_ART/prisma/miti_art"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandle(c *gin.Context, prisma *miti_art.PrismaClient) {
	// RegisterRequest struct for client registration
	var req struct {
		ClientFirstName string `json:"clientFirstName" binding:"required"`
		ClientOtherName string `json:"clientOtherName" binding:"required"`
		ClientEmail     string `json:"clientEmail" binding:"required,email"`
		ClientPassword  string `json:"clientPassword" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	service.RegisterClient(prisma, req.ClientEmail, req.ClientFirstName, req.ClientOtherName, req.ClientPassword)

}
