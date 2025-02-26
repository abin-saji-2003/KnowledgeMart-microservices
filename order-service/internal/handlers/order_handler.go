package handlers

import (
	"context"
	"log"
	"order-service/internal/usecase"

	orderProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/order-pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandler struct {
	orderUC *usecase.OrderUseCase
	orderProto.UnimplementedOrderServiceServer
}

func NewOrderHandler(orderUC *usecase.OrderUseCase) *OrderHandler { // ✅ Expect pointer type
	return &OrderHandler{orderUC: orderUC}
}

func (h *OrderHandler) PlaceOrder(ctx context.Context, req *orderProto.PlaceOrderRequest) (*orderProto.PlaceOrderResponse, error) {
	log.Printf("Received PlaceOrder request for User ID: %d", req.UserId)

	order, orderItems, err := h.orderUC.PlaceOrder(uint(req.UserId))
	if err != nil {
		return &orderProto.PlaceOrderResponse{
			Status:  "Failed",
			Message: err.Error(),
			Order:   nil,
		}, err
	}

	var protoOrderItems []*orderProto.OrderItem
	for _, item := range *orderItems {
		protoOrderItems = append(protoOrderItems, &orderProto.OrderItem{
			OrderId:     uint32(item.OrderID),
			ProductId:   uint32(item.ProductID),
			Price:       float32(item.Price),
			FinalAmount: float32(item.FinalAmount),
			Status:      item.Status,
		})
	}

	// Convert order model to protobuf format
	orderProtoData := &orderProto.Order{
		Id:             uint32(order.OrderID),
		UserId:         uint32(order.UserID),
		Items:          protoOrderItems,
		TotalPrice:     float32(order.TotalAmount),
		FinalAmount:    float32(order.FinalAmount),
		PaymentMethod:  order.PaymentMethod,
		PaymentStatus:  order.PaymentStatus,
		OrderedAt:      timestamppb.New(order.OrderedAt),
		Status:         order.Status,
		DeliveryCharge: float32(order.DeliveryCharge),
	}

	return &orderProto.PlaceOrderResponse{
		Status:  "Success",
		Message: "Order placed successfully",
		Order:   orderProtoData,
	}, nil
}

func (h *OrderHandler) GetOrders(ctx context.Context, req *orderProto.GetOrdersRequest) (*orderProto.GetOrdersResponse, error) {
	orders, err := h.orderUC.GetUserOrders(uint(req.UserId))
	if err != nil {
		return &orderProto.GetOrdersResponse{
			Status:  "failed",
			Message: err.Error(),
		}, err
	}

	var orderList []*orderProto.Order
	for _, order := range orders {
		var items []*orderProto.OrderItem
		for _, item := range order.OrderItems {
			items = append(items, &orderProto.OrderItem{
				OrderId:     uint32(item.OrderID),
				ProductId:   uint32(item.ProductID),
				SellerId:    uint32(item.SellerID),
				Name:        item.Name,
				Description: item.Description,
				Price:       float32(item.Price),
				FinalAmount: float32(item.FinalAmount),
				Status:      item.Status,
				ImageUrl:    item.Image,
			})
		}

		orderList = append(orderList, &orderProto.Order{
			Id:             uint32(order.OrderID),
			UserId:         uint32(order.UserID),
			Items:          items,
			TotalPrice:     float32(order.TotalAmount),
			FinalAmount:    float32(order.FinalAmount),
			PaymentMethod:  order.PaymentMethod,
			PaymentStatus:  order.PaymentStatus,
			OrderedAt:      timestamppb.New(order.OrderedAt),
			Status:         order.Status,
			DeliveryCharge: float32(order.DeliveryCharge),
		})
	}

	log.Print("check", orderList)
	return &orderProto.GetOrdersResponse{
		Status:  "success",
		Message: "Orders retrieved successfully",
		Orders:  orderList,
	}, nil
}
