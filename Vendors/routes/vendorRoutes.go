package routes

import (
	controller "MITI_ART/Vendors/controller"
	"MITI_ART/prisma/miti_art"

	"github.com/gin-gonic/gin"
)

func VendorsRoutes(router *gin.Engine, prisma *miti_art.PrismaClient) {
	user := router.Group("/vendor")
	{
		user.POST("/register", func(c *gin.Context) {
			controller.RegisterHandle(c, prisma)
		})
	}
}
