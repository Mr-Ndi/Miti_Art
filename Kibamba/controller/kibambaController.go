package controllers

import (
	"context"
	"net/http"

	"MITI_ART/prisma/db"
	"MITI_ART/services"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context, prisma *db.PrismaClient) {
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
