package main

import (
	"fmt"

	"MITI_ART/prisma/db"
	"MITI_ART/routes"
	"MITI_ART/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	prisma := db.NewClient()
	defer prisma.Disconnect()

	// Seed admin user
	services.SeedAdmin(prisma)

	// Setup routes
	routes.AdminRoutes(r, prisma)

	fmt.Println("Server is running on port 8080")
	r.Run(":8080")
}
