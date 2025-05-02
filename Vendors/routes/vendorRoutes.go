package routes

import (
	controller "MITI_ART/Vendors/controller"
	middlewares "MITI_ART/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VendorsRoutes(router *gin.Engine, db *gorm.DB) {
	vendor := router.Group("/vendor")
	{
		vendor.POST("/register", func(c *gin.Context) {
			controller.RegisterHandle(c, db)
		})
		auth := vendor.Group("", middlewares.AuthMiddleware())
		{
			//Router for enabling vendor to upload a product
			auth.POST("/upload", func(c *gin.Context) { controller.UploadHandle(c, db) })
			//Router for getting the product that belongs to the vendor
			auth.GET("/my-products", func(c *gin.Context) {})
			//Router for retriving a single product
			auth.POST("/my-product/:id", func(c *gin.Context) {})
			//Router for retriving all orders recomended from his product
			auth.GET("/required-product", func(c *gin.Context) {})
			//Router for editing the product description
			auth.POST("/edit-product/:id", func(c *gin.Context) {})
			//Router for deleting the posted product
			auth.POST("/remove-product/:id", func(c *gin.Context) {})
		}
	}
}
