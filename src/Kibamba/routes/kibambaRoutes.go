package routes

import (
	controllers "MITI_ART/src/Kibamba/controller"
	middlewares "MITI_ART/src/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(router *gin.Engine, db *gorm.DB) {
	admin := router.Group("/admin")
	// admin.Use(middlewares.AuthMiddleware())
	{
		admin.POST("/invite", controllers.InvitationHandler)
		kibamba := admin.Group("", middlewares.RequireAdmin())
		{
			kibamba.GET("/view-clients", func(c *gin.Context) { controllers.ViewClients(c, db) })
			kibamba.GET("/view-vendors", func(c *gin.Context) { controllers.ViewVendors(c, db) })
			kibamba.GET("/view-orders", func(c *gin.Context) { controllers.ViewOrders(c, db) })
			kibamba.GET("/view-products", func(c *gin.Context) { controllers.ViewAllProducts(c, db) })
			kibamba.POST("/edit-vendor", func(c *gin.Context) { controllers.EditVendor(c, db) })
			kibamba.POST("/edit-client", func(c *gin.Context) { controllers.EditClient(c, db) })
			kibamba.POST("/eliminate-vendor", func(c *gin.Context) { controllers.EliminateVendor(c, db) })
			kibamba.POST("/eliminate-client", func(c *gin.Context) { controllers.EliminateClient(c, db) })
		}
	}
	user := router.Group("/user")
	{
		user.POST("/login", func(c *gin.Context) {
			controllers.LoginHandler(c, db)
		})
	}
}
