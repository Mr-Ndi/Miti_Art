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
		kibamba := admin.Group("", middlewares.AuthMiddleware())
		{
			kibamba.GET("/view-clients", func(c *gin.Context) { controllers.ViewClients(c, db) })
			kibamba.GET("/view-vendors", func(c *gin.Context) { controllers.ViewVendors(c, db) })
			kibamba.GET("/view-orders", func(c *gin.Context) {})
			kibamba.GET("/view-products", func(c *gin.Context) {})
			kibamba.POST("/edit-vendor", func(c *gin.Context) {})
			kibamba.POST("/edit-client", func(c *gin.Context) {})
			kibamba.POST("/eliminate-vendor", func(c *gin.Context) {})
			kibamba.POST("/eliminate-client", func(c *gin.Context) {})
		}
	}
	user := router.Group("/user")
	{
		user.POST("/login", func(c *gin.Context) {
			controllers.LoginHandler(c, db)
		})
	}
}
