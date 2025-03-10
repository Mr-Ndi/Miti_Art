package routes

import (
	controllers "MITI_ART/Kibamba/controller"
	"MITI_ART/prisma/miti_art"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine, prisma *miti_art.PrismaClient) {
	admin := router.Group("/admin")
	{
		admin.POST("/login", func(c *gin.Context) {
			controllers.LoginHandler(c, prisma)
		})
		admin.POST("/invite", controllers.InvitationHandler)
	}
}
