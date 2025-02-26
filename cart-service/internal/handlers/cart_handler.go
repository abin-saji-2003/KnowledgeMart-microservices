package handlers

import (
	"cart-service/internal/usecase"
	"context"

	cartProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/cart-pb"
)

type CartHandlers struct {
	cartUC *usecase.CartUseCase // ✅ Use a pointer to CartUseCase struct
	cartProto.UnimplementedCartServiceServer
}

func NewCartHandler(cartUC *usecase.CartUseCase) *CartHandlers { // ✅ Accepts a pointer
	return &CartHandlers{cartUC: cartUC}
}

func (h *CartHandlers) AddToCart(ctx context.Context, req *cartProto.AddToCartRequest) (*cartProto.CartResponse, error) {
	err := h.cartUC.AddToCart(uint(req.UserId), uint(req.ProductId))
	if err != nil {
		return &cartProto.CartResponse{Status: "failed", Message: err.Error()}, err
	}

	return &cartProto.CartResponse{
		Status:  "success",
		Message: "Product added to cart successfully",
	}, nil
}

func (h *CartHandlers) GetProductsFromCart(ctx context.Context, req *cartProto.GetProductsFromCartRequest) (*cartProto.GetProductsFromCartResponse, error) {
	products, err := h.cartUC.GetProductsFromCart(uint(req.UserId))
	if err != nil {
		return &cartProto.GetProductsFromCartResponse{
			Status:   "error",
			Message:  "Failed to fetch products from cart",
			Products: nil,
		}, err
	}

	var cartList []*cartProto.CartProduct
	for _, c := range products {
		cartList = append(cartList, &cartProto.CartProduct{
			ProductId:    uint32(c.ProductID),
			Name:         c.Name,
			Description:  c.Description,
			Price:        float32(c.Price),
			OfferAmount:  float32(c.OfferAmount),
			ImageUrl:     c.ImageURL,
			Availability: c.Availability,
		})
	}

	return &cartProto.GetProductsFromCartResponse{
		Status:   "success",
		Message:  "Cart products retrieved successfully",
		Products: cartList,
	}, nil
}

func (h *CartHandlers) RemoveFromCart(ctx context.Context, req *cartProto.RemoveFromCartRequest) (*cartProto.CartResponse, error) {
	err := h.cartUC.RemoveFromCart(uint(req.ProductId), uint(req.UserId))
	if err != nil {
		return &cartProto.CartResponse{
			Status:  "failed",
			Message: err.Error(),
		}, err
	}
	return &cartProto.CartResponse{
		Status:  "success",
		Message: "Product successfully removed from the cart",
	}, nil
}

func (h *CartHandlers) ClearCart(ctx context.Context, req *cartProto.ClearCartRequest) (*cartProto.CartResponse, error) {
	err := h.cartUC.ClearCart(uint(req.UserId))
	if err != nil {
		return &cartProto.CartResponse{Status: "error", Message: err.Error()}, nil
	}
	return &cartProto.CartResponse{Status: "success", Message: "Cart cleared"}, nil
}
