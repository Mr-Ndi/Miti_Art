package routes

import (
	"MITI_ART/controllers"
	"MITI_ART/prisma/db"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine, prisma *db.PrismaClient) {
	router.POST("/login", func(c *gin.Context) { controllers.LoginHandler(c, prisma) })
}
