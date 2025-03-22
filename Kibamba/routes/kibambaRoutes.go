package routes

import (
	controllers "MITI_ART/Kibamba/controller"
	middlewares "MITI_ART/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(router *gin.Engine, db *gorm.DB) {
	admin := router.Group("/admin")
	admin.Use(middlewares.AuthMiddleware())
	{
		admin.POST("/invite", controllers.InvitationHandler)
	}
	user := router.Group("/user")
	{
		user.POST("/login", func(c *gin.Context) {
			controllers.LoginHandler(c, db)
		})
	}
}
