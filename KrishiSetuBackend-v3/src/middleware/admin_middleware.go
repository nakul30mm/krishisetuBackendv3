package middleware

import (
	"krishisetu-backend/database"
	"krishisetu-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminCheck ensures the user has administrative privileges
func AdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// BlockedCheck ensures the user's account is not suspended
func BlockedCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next() // Should be caught by JWTAuth if required
			return
		}

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.Next()
			return
		}

		if user.IsBlocked {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account suspended. Please contact support."})
			c.Abort()
			return
		}

		c.Next()
	}
}
