package handlers

import (
	"api-gateway/internal/middleware"
	"context"
	"net/http"
	"strconv"

	cartProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/cart-pb"
	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(router *gin.Engine, cartClient cartProto.CartServiceClient) {
	cartRoutes := router.Group("/cart")
	{
		cartRoutes.POST("/add", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			// Extract user ID from JWT middleware
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
				return
			}

			// Parse request body
			var req cartProto.AddToCartRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			// Assign the extracted user ID
			req.UserId = uint32(userID.(uint))

			// Call CartService via gRPC
			resp, err := cartClient.AddToCart(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Return response
			c.JSON(http.StatusOK, gin.H{
				"status":  resp.Status,
				"message": resp.Message,
			})
		})

		cartRoutes.GET("/", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			// Extract user ID from JWT middleware
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
				return
			}

			// Call gRPC CartService to get cart items
			resp, err := cartClient.GetProductsFromCart(context.Background(), &cartProto.GetProductsFromCartRequest{
				UserId: uint32(userID.(uint)),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Return response
			c.JSON(http.StatusOK, gin.H{
				"status":   resp.Status,
				"message":  resp.Message,
				"products": resp.Products,
			})
		})

		cartRoutes.DELETE("/delete/:id", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
				return
			}

			productID, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
				return
			}

			req := &cartProto.RemoveFromCartRequest{
				ProductId: uint32(productID),
				UserId:    uint32(userID.(uint)),
			}

			resp, err := cartClient.RemoveFromCart(context.Background(), req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Send response
			c.JSON(http.StatusOK, gin.H{
				"status":  resp.Status,
				"message": resp.Message,
			})
		})
	}
}
