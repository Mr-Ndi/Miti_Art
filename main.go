package main

import (
	"MITI_ART/Kibamba/routes"
	database "MITI_ART/configure"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{})

	// database stuffs
	prisma := database.InitDB()
	defer prisma.Disconnect()

	// Setup routes
	routes.AdminRoutes(r, prisma)

	fmt.Println("Server is running on port 8080")
	r.Run(":8080")
}
