package main

import (
	"fmt"
	"os"

	database "MITI_ART/configure"
	route "MITI_ART/src/routes"

	_ "MITI_ART/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title My API
// @version 1.0
// @description This is my API using Gin, GORM, and Swagger
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your token.

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

	// Connect to DB
	database.ConnectDB()

	// App Routes
	route.AllRoutes(r, database.DB)

	// Swagger docs route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("Connected to PostgreSQL and migrated successfully!")
	fmt.Println("--------------------------------------------------------------")
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on: http://%s:%s\n", host, port)
	fmt.Printf("Swagger docs available at: http://%s:%s/swagger/index.html\n", host, port)

	r.Run(":8080")
}
