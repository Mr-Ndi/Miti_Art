package routes

import (
	controller "MITI_ART/Client/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ClientRoutes(router *gin.Engine, db *gorm.DB) {
	user := router.Group("/user")
	{
		user.POST("/register", func(c *gin.Context) {
			controller.RegisterHandle(c, db)
		})
		user.GET("/furniture", func(c *gin.Context) {
			controller.RegisterHandle(c, db)
		})
	}
}
