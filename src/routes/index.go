package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	client "MITI_ART/src/Client/routes"
	kibamba "MITI_ART/src/Kibamba/routes"
	vendor "MITI_ART/src/Vendors/routes"
)

func AllRoutes(r *gin.Engine, db *gorm.DB) {
	client.ClientRoutes(r, db)
	kibamba.AdminRoutes(r, db)
	vendor.VendorsRoutes(r, db)
}
