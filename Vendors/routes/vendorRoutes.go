package routes

import (
	controller "MITI_ART/Vendors/controller"
	"MITI_ART/prisma/miti_art"

	"github.com/gin-gonic/gin"
)

func VendorsRoutes(router *gin.Engine, prisma *miti_art.PrismaClient) {
	vendor := router.Group("/vendor")
	{
		vendor.POST("/register", func(c *gin.Context) {
			controller.RegisterHandle(c, prisma)
		})
	}
}
