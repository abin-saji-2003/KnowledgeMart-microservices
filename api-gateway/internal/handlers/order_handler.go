package handlers

import (
	"api-gateway/internal/middleware"
	"context"
	"net/http"

	orderProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/order-pb"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderClient orderProto.OrderServiceClient
}

// NewOrderHandler initializes the order handler
func NewOrderHandler(orderClient orderProto.OrderServiceClient) *OrderHandler {
	return &OrderHandler{orderClient: orderClient}
}

// RegisterOrderRoutes sets up order routes
func (h *OrderHandler) RegisterOrderRoutes(r *gin.Engine) {
	orderRoutes := r.Group("/api/order")
	orderRoutes.Use(middleware.JWTAuthMiddleware()) // Protect routes
	{
		orderRoutes.POST("/place", h.PlaceOrder)
		orderRoutes.GET("/list", h.GetOrders)
	}
}

// PlaceOrder handles order placement
func (h *OrderHandler) PlaceOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID missing from token"})
		return
	}

	// Call gRPC order service
	resp, err := h.orderClient.PlaceOrder(context.Background(), &orderProto.PlaceOrderRequest{
		UserId: uint32(userID.(uint)),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": resp.Status, "message": resp.Message, "order": resp.Order})
}

// GetOrders retrieves user orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID missing from token"})
		return
	}

	// Call gRPC order service
	resp, err := h.orderClient.GetOrders(context.Background(), &orderProto.GetOrdersRequest{
		UserId: uint32(userID.(uint)),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": resp.Status, "message": resp.Message, "orders": resp.Orders})
}
