package routes

import (
	"your_project/controllers"
	"your_project/prisma/db"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine, prisma *db.PrismaClient) {
	router.POST("/login", func(c *gin.Context) { controllers.LoginHandler(c, prisma) })
}
