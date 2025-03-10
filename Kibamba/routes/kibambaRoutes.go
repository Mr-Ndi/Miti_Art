package routes

import (
	controllers "MITI_ART/Kibamba/controller"
	middlewares "MITI_ART/middleware"
	"MITI_ART/prisma/miti_art"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine, prisma *miti_art.PrismaClient) {
	admin := router.Group("/admin")
	admin.Use(middlewares.AuthMiddleware())
	{
		admin.POST("/invite", controllers.InvitationHandler)
	}
	user := router.Group("/user")
	{
		user.POST("/login", func(c *gin.Context) {
			controllers.LoginHandler(c, prisma)
		})
	}
}
