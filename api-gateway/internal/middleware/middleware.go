package middleware

import (
	"api-gateway/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token format",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Token validation failed: %v", err),
			})
			c.Abort()
			return
		}

		switch claims.Role {
		case "user":
			c.Set("userID", claims.ID)
			c.Set("role", "user")
		case "seller":
			c.Set("sellerID", claims.ID)
			c.Set("role", "seller")
		default:
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Unauthorized role",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
