package middlewares

import (
	"net/http"
	"os"
	"strings"

	utils "MITI_ART/Utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		email, emailExists := claims["email"].(string)
		role, roleExists := claims["role"].(string)
		userIDStr, idExists := claims["user_id"].(string)

		if !emailExists || !roleExists || !idExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		c.Set("userEmail", email)
		c.Set("userRole", role)
		c.Set("user_id", userID)

		c.Next()
	}
}
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail, exists := c.Get("user_email")
		if !exists || userEmail != os.Getenv("ADMIN_EMAIL") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access only"})
			return
		}
		c.Next()
	}
}
