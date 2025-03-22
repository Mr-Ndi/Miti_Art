package main

import (
	client "MITI_ART/Client/routes"
	kibamba "MITI_ART/Kibamba/routes"
	vendor "MITI_ART/Vendors/routes"
	database "MITI_ART/configure"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{})

	// database stuffs
	database.ConnectDB()

	// Routes setup
	kibamba.AdminRoutes(r, database.DB)
	client.ClientRoutes(r, database.DB)
	vendor.VendorsRoutes(r, database.DB)

	fmt.Println("Server is running on port 8080")
	r.Run(":8080")
}
