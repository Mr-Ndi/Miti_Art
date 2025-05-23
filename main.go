package main

import (
	client "MITI_ART/Client/routes"
	kibamba "MITI_ART/Kibamba/routes"
	vendor "MITI_ART/Vendors/routes"
	database "MITI_ART/configure"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{})
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}
	r.Use(cors.New(corsConfig))

	// database stuffs
	database.ConnectDB()

	// Routes setup
	kibamba.AdminRoutes(r, database.DB)
	client.ClientRoutes(r, database.DB)
	vendor.VendorsRoutes(r, database.DB)

	fmt.Println("Server is running on port 8080")
	r.Run(":8080")
}
