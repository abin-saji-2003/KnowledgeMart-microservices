package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"order-service/internal/models"
	"order-service/internal/repository"

	cartProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/cart-pb"
	productProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/product-pb"
)

type OrderUseCase struct {
	orderRepo     repository.OrderRepository
	cartClient    cartProto.CartServiceClient
	productClient productProto.ProductServiceClient
}

func NewOrderUseCase(orderRepo repository.OrderRepository, cartClient cartProto.CartServiceClient, productClient productProto.ProductServiceClient) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:     orderRepo,
		cartClient:    cartClient,
		productClient: productClient,
	}
}
func (uc *OrderUseCase) PlaceOrder(userID uint) (*models.Order, *[]models.OrderItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cartRes, err := uc.cartClient.GetProductsFromCart(ctx, &cartProto.GetProductsFromCartRequest{UserId: uint32(userID)})
	if err != nil {
		return nil, nil, errors.New("failed to fetch cart items: " + err.Error())
	}
	if len(cartRes.Products) == 0 {
		return nil, nil, errors.New("cart is empty")
	}

	var totalAmount float64
	for _, item := range cartRes.Products {
		if !item.Availability {
			return nil, nil, errors.New("some items in the cart are out of stock")
		}
		totalAmount += float64(item.Price)
	}

	finalAmount := totalAmount
	deliveryCharge := 0
	if finalAmount < 500 {
		deliveryCharge = 40
		finalAmount += float64(deliveryCharge)
	}

	order := models.Order{
		UserID:         userID,
		TotalAmount:    totalAmount,
		FinalAmount:    finalAmount,
		PaymentMethod:  "COD",
		PaymentStatus:  "Pending",
		OrderedAt:      time.Now(),
		Status:         "confirmed",
		DeliveryCharge: float64(deliveryCharge),
	}

	orderID, err := uc.orderRepo.CreateOrder(&order)
	if err != nil {
		return nil, nil, err
	}

	var orderItems []models.OrderItem
	var firstSellerID uint

	for _, cartItem := range cartRes.Products {
		productRes, err := uc.productClient.GetProductById(ctx, &productProto.GetProductRequest{ProductId: uint32(cartItem.ProductId)})
		if err != nil || productRes.Data == nil {
			return nil, nil, errors.New("failed to fetch product details for product ID: " + fmt.Sprint(cartItem.ProductId))
		}
		if firstSellerID == 0 {
			firstSellerID = uint(cartItem.SellerId)
		}

		orderItem := models.OrderItem{
			OrderID:     orderID,
			ProductID:   uint(cartItem.ProductId),
			UserID:      userID,
			Price:       float64(cartItem.Price),
			FinalAmount: float64(cartItem.Price),
			Status:      "confirmed",
			SellerID:    uint(cartItem.SellerId),
		}

		err = uc.orderRepo.CreateOrderItems(orderItem)
		if err != nil {
			return nil, nil, err
		}

		_, err = uc.productClient.EditProduct(ctx, &productProto.EditProductRequest{
			ProductId:    cartItem.ProductId,
			Availability: &[]bool{false}[0],
			SellerId:     14,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update product availability for product ID: %d, error: %s", cartItem.ProductId, err.Error())
		}

		orderItems = append(orderItems, orderItem)
	}

	order.SellerID = firstSellerID
	err = uc.orderRepo.UpdateOrderSellerID(orderID, firstSellerID)
	if err != nil {
		return nil, nil, err
	}

	_, err = uc.cartClient.ClearCart(ctx, &cartProto.ClearCartRequest{UserId: uint32(userID)})
	if err != nil {
		return nil, nil, errors.New("failed to clear cart: " + err.Error())
	}

	return &order, &orderItems, nil
}

func (uc *OrderUseCase) GetUserOrders(userID uint) ([]models.OrderResponse, error) {
	orders, err := uc.orderRepo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve user orders: " + err.Error())
	}

	var orderResponses []models.OrderResponse

	for _, order := range orders {
		orderItems, err := uc.orderRepo.GetOrderItemsByOrderID(order.OrderID)
		if err != nil {
			return nil, errors.New("failed to retrieve order items: " + err.Error())
		}

		var itemResponses []models.OrderItemResponse
		for _, item := range orderItems {
			productRes, err := uc.productClient.GetProductById(context.Background(), &productProto.GetProductRequest{ProductId: uint32(item.ProductID)})
			if err != nil || productRes.Data == nil {
				return nil, fmt.Errorf("failed to fetch product details for product ID: %d", item.ProductID)
			}

			itemResponses = append(itemResponses, models.OrderItemResponse{
				OrderID:     item.OrderID,
				UserID:      item.UserID,
				ProductID:   item.ProductID,
				SellerID:    item.SellerID,
				Name:        productRes.Data.Name,
				Description: productRes.Data.Description,
				Price:       item.Price,
				FinalAmount: item.FinalAmount,
				Status:      item.Status,
				Image:       productRes.Data.ImageUrl,
			})
		}

		orderResponses = append(orderResponses, models.OrderResponse{
			OrderID:        order.OrderID,
			UserID:         order.UserID,
			DeliveryCharge: order.DeliveryCharge,
			TotalAmount:    order.TotalAmount,
			FinalAmount:    order.FinalAmount,
			PaymentMethod:  order.PaymentMethod,
			PaymentStatus:  order.PaymentStatus,
			OrderedAt:      order.OrderedAt,
			SellerID:       order.SellerID,
			Status:         order.Status,
			OrderItems:     itemResponses,
		})
	}

	return orderResponses, nil
}
