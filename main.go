package main

import (
	client "MITI_ART/Client/routes"
	kibamba "MITI_ART/Kibamba/routes"
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
	kibamba.AdminRoutes(r, prisma)
	client.ClientRoutes(r, prisma)

	fmt.Println("Server is running on port 8080")
	r.Run(":8080")
}
