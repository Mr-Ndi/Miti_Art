package routes

import (
	controller "MITI_ART/src/Vendors/controller"
	middlewares "MITI_ART/src/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VendorsRoutes(router *gin.Engine, db *gorm.DB) {
	vendor := router.Group("/vendor")
	{
		//Public routes
		vendor.POST("/register", func(c *gin.Context) {
			controller.RegisterHandle(c, db)
		})
		//Routes available for Authenticated users
		auth := vendor.Group("", middlewares.AuthMiddleware())
		{
			//Router for enabling vendor to upload a product
			auth.POST("/upload", func(c *gin.Context) { controller.UploadHandle(c, db) })
			//Router for getting the product that belongs to the vendor
			auth.GET("/my-products", func(c *gin.Context) { controller.MyProducts(c, db) })
			//Router for retriving a single product
			auth.GET("/my-product/:id", func(c *gin.Context) { controller.MyProduct(c, db) })
			//Router for retriving all orders recomended from his product
			auth.GET("/required-product", func(c *gin.Context) { controller.MyOrders(c, db) })
			//Router for editing the product description
			auth.PATCH("/edit-product/:id", func(c *gin.Context) { controller.EditProduct(c, db) })
			//Router for deleting the posted product
			auth.DELETE("/remove-product/:id", func(c *gin.Context) { controller.DeleteProduct(c, db) })
		}
	}
}
