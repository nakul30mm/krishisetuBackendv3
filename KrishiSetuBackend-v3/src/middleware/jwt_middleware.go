package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"krishisetu-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// validates and verifies the jwt token on the header of user request
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JWTSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		rawID, ok := claims["user_id"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id missing in token"})
			c.Abort()
			return
		}

		var userID uint
		switch v := rawID.(type) {
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		case string:
			fmt.Sscan(v, &userID)
		default:
			userID = 0
		}

		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", claims["email"])

		c.Next()
	}
}
