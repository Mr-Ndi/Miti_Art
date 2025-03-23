package routes

import (
	controller "MITI_ART/Vendors/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VendorsRoutes(router *gin.Engine, db *gorm.DB) {
	vendor := router.Group("/vendor")
	{
		vendor.POST("/register", func(c *gin.Context) {
			controller.RegisterHandle(c, db)
		})
		vendor.POST("/upload", func(c *gin.Context) {
			controller.UploadHandle(c, db)
		})
	}
}
