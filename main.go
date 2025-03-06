package main

import (
	"fmt"

	"MITI_ART/Kibamba/routes"
	"MITI_ART/prisma/miti_art"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{})
	prisma := miti_art.NewClient()
	defer prisma.Disconnect()

	// Setup routes
	routes.AdminRoutes(r, prisma)

	fmt.Println("Server is running on port 8080")
	r.Run(":8080")
}
