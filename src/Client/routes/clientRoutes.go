package routes

import (
	controller "MITI_ART/src/Client/controller"
	middleware "MITI_ART/src/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ClientRoutes(router *gin.Engine, db *gorm.DB) {
	user := router.Group("/user")
	{
		// Unprotected route!
		user.POST("/register", func(c *gin.Context) {
			controller.RegisterHandle(c, db)
		})
		// Unprotected route!
		user.GET("/furniture", func(c *gin.Context) {
			controller.GetFurniture(c, db)
		})
		// Unprotected route!
		user.GET("/furniture/:id", func(c *gin.Context) {
			controller.GetFurnitureDetails(c, db)
		})
		// Secured the routes!
		auth := user.Group("", middleware.AuthMiddleware())
		{
			auth.POST("/order/:id", func(c *gin.Context) { controller.CreateOrder(c, db) })
			auth.GET("/wished-item/:id", func(c *gin.Context) { controller.AppendWishList(c, db) })
			auth.GET("/my-orders", func(c *gin.Context) { controller.ListUserOrders(c, db) })

		}
	}
}
