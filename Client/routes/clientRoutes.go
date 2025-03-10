package routes

import (
	controller "MITI_ART/Client/controller"
	"MITI_ART/prisma/miti_art"

	"github.com/gin-gonic/gin"
)

func ClientRoutes(router *gin.Engine, prisma *miti_art.PrismaClient) {
	user := router.Group("/user")
	{
		user.POST("/register", func(c *gin.Context) {
			controller.RegisterHandle(c, prisma)
		})
	}
}
