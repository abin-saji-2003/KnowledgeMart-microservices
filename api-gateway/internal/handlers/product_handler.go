package handlers

import (
	"api-gateway/internal/middleware"
	"context"
	"net/http"
	"strconv"

	productProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/product-pb"

	"github.com/gin-gonic/gin"
)

func AddProduct(r *gin.Engine, productClient productProto.ProductServiceClient) {
	productRoutes := r.Group("/api/product")
	{
		productRoutes.POST("/add", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			// Extract sellerID from JWT middleware
			sellerID, exists := c.Get("sellerID")
			if !exists {
				c.JSON(http.StatusForbidden, gin.H{"error": "Only sellers can add products"})
				return
			}

			var req productProto.AddProductRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			// Assign the correct seller ID
			req.SellerId = uint32(sellerID.(uint))

			// // Debugging Log
			// c.JSON(http.StatusOK, gin.H{"debug_request": req})

			resp, err := productClient.AddProduct(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Send back the correct response
			c.JSON(http.StatusOK, gin.H{
				"status":  resp.Status,
				"message": resp.Message,
				"product": gin.H{
					"id":           resp.Data.Id,
					"name":         resp.Data.Name,
					"description":  resp.Data.Description,
					"availability": resp.Data.Availability,
					"price":        resp.Data.Price,
					"offer_amount": resp.Data.OfferAmount,
					"image_url":    resp.Data.ImageUrl,
				},
			})
		})

		productRoutes.GET("/all", func(c *gin.Context) {
			resp, err := productClient.GetAllProducts(context.Background(), &productProto.Empty{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":   resp.Status,
				"message":  resp.Message,
				"products": resp.Products,
			})
		})

		productRoutes.PUT("/edit", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			// Extract sellerID from JWT middleware
			sellerID, exists := c.Get("sellerID")
			if !exists {
				c.JSON(http.StatusForbidden, gin.H{"error": "Only sellers can edit products"})
				return
			}

			var req productProto.EditProductRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			// Assign the correct seller ID
			req.SellerId = uint32(sellerID.(uint))

			// Call gRPC service to edit product
			resp, err := productClient.EditProduct(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Send back the correct response
			c.JSON(http.StatusOK, gin.H{
				"status":  resp.Status,
				"message": resp.Message,
				"product": gin.H{
					"id":           resp.Data.Id,
					"name":         resp.Data.Name,
					"description":  resp.Data.Description,
					"availability": resp.Data.Availability,
					"price":        resp.Data.Price,
					"offer_amount": resp.Data.OfferAmount,
					"image_url":    resp.Data.ImageUrl,
				},
			})
		})

		productRoutes.DELETE("/delete/:id", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			// Extract sellerID from JWT middleware
			sellerID, exists := c.Get("sellerID")
			if !exists {
				c.JSON(http.StatusForbidden, gin.H{"error": "Only sellers can delete products"})
				return
			}

			// Get Product ID from URL parameter
			productID, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
				return
			}

			// Prepare gRPC request
			req := &productProto.DeleteProductRequest{
				ProductId: uint32(productID),
				SellerId:  uint32(sellerID.(uint)),
			}

			// Call gRPC service
			resp, err := productClient.DeleteProduct(context.Background(), req)
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
