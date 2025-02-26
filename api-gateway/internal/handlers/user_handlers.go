package handlers

import (
	"api-gateway/internal/middleware"
	"context"
	"net/http"

	authProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/auth-pb"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up authentication-related API routes
func RegisterAuthRoutes(r *gin.Engine, authClient authProto.AuthServiceClient) {
	authRoutes := r.Group("/api/auth")
	{
		// User Signup
		authRoutes.POST("/signup", func(c *gin.Context) {
			var req authProto.EmailSignupRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			resp, err := authClient.EmailSignup(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": resp.Message, "success": resp.Success})
		})

		// User Login
		authRoutes.POST("/login", func(c *gin.Context) {
			var req authProto.EmailLoginRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			resp, err := authClient.EmailLogin(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"token":   resp.Token,
				"success": resp.Success,
				"user": gin.H{
					"id":           resp.User.Id,
					"name":         resp.User.Name,
					"email":        resp.User.Email,
					"phone_number": resp.User.PhoneNumber,
					"picture":      resp.User.Picture,
					"blocked":      resp.User.Blocked,
					"verified":     resp.User.Verified,
				},
			})
		})

		// Seller Signup (Protected)
		authRoutes.POST("/seller/signup", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			var req authProto.SellerRegisterRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			// Extract UserID from JWT middleware
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			// Attach extracted UserID to gRPC request
			req.UserId = uint32(userID.(uint))

			resp, err := authClient.SellerRegister(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": resp.Message,
				"success": resp.Success,
			})
		})

		// Seller Login
		authRoutes.POST("/seller/login", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			var req authProto.SellerLoginRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			req.UserId = uint32(userID.(uint))

			resp, err := authClient.SellerLogin(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token":   resp.Token,
				"success": resp.Success,
				"seller": gin.H{
					"id":          resp.Seller.Id,
					"name":        resp.Seller.Name,
					"description": resp.Seller.Description,
					"user_id":     resp.Seller.UserId,
				},
			})
		})
	}
}
