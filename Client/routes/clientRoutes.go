package routes

import (
	controller "MITI_ART/Client/controller"
	middleware "MITI_ART/middleware"

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
		// Secure the order route!
		user.POST("/order/:id", middleware.AuthMiddleware(), func(c *gin.Context) {
			controller.CreateOrder(c, db)
		})
		// Secure the wishlist route!
		user.POST("/wished-item/:id", middleware.AuthMiddleware(), func(c *gin.Context) {
			controller.AppendWishList(c, db)
		})
		// Secure the Orders!
		user.POST("/my-orders", middleware.AuthMiddleware(), func(c *gin.Context) {
			controller.GetOrders(c, db)
		})
	}
}
