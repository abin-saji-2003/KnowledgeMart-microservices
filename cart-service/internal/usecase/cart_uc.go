package usecase

import (
	"cart-service/internal/models"
	"cart-service/internal/repository"
	"context"
	"errors"
	"log"

	productProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/product-pb"
)

type CartUseCase struct {
	cartRepo      repository.CartRepository
	productClient productProto.ProductServiceClient
}

func NewCartUsecase(cartRepo repository.CartRepository, productClient productProto.ProductServiceClient) *CartUseCase {
	return &CartUseCase{cartRepo: cartRepo, productClient: productClient}
}

func (cu *CartUseCase) AddToCart(userId, productId uint) error {
	log.Printf("Calling Product Service with product ID: %d", productId)

	resp, err := cu.productClient.GetProductById(context.Background(), &productProto.GetProductRequest{ProductId: uint32(productId)})
	if err != nil {
		log.Printf("gRPC error: %v", err) // Log full gRPC error details
		return errors.New("failed to fetch product details from product service")
	}

	// Check if product exists and is available
	if resp.Data == nil {
		log.Println("Product not found in Product Service")
		return errors.New("product not found")
	}

	if !resp.Data.Availability {
		log.Println("Product is not available")
		return errors.New("product is not available")
	}

	exists, err := cu.cartRepo.IsProductInCart(userId, productId)
	if err != nil {
		log.Println("Error checking cart:", err)
		return errors.New("failed to check cart")
	}

	if exists {
		log.Println("Product already exists in cart")
		return errors.New("product already exists in cart")
	}

	cart := models.Cart{
		UserID:    userId,
		ProductID: productId,
	}

	// Add product to the cart
	if err := cu.cartRepo.AddtoCart(&cart); err != nil {
		log.Println("Failed to add product to cart")
		return errors.New("failed to add product to cart")
	}

	log.Println("Product successfully added to cart")
	return nil
}

func (cu *CartUseCase) GetProductsFromCart(userID uint) ([]models.CartProduct, error) {
	// Get cart items (only ProductIDs)
	cartItems, err := cu.cartRepo.GetProductsFromCart(userID)
	if err != nil {
		return nil, err
	}

	var cartProducts []models.CartProduct
	for _, cart := range cartItems {
		// Fetch product details from the Product Service using gRPC
		resp, err := cu.productClient.GetProductById(context.Background(), &productProto.GetProductRequest{ProductId: uint32(cart.ProductID)})
		if err != nil || resp.Data == nil {
			log.Printf("Failed to fetch product details for ProductID %d: %v", cart.ProductID, err)
			continue // Skip this product but continue with others
		}

		// Append full product details
		cartProducts = append(cartProducts, models.CartProduct{
			ProductID:    cart.ProductID,
			Name:         resp.Data.Name,
			Description:  resp.Data.Description,
			Price:        float64(resp.Data.Price),
			OfferAmount:  float64(resp.Data.OfferAmount),
			ImageURL:     resp.Data.ImageUrl,
			Availability: resp.Data.Availability,
		})
	}

	return cartProducts, nil
}

func (cu *CartUseCase) RemoveFromCart(productID, userID uint) error {
	err := cu.cartRepo.RemoveFromCart(productID, userID)
	if err != nil {
		return errors.New("failed to remove product")
	}
	return nil
}

func (uc *CartUseCase) ClearCart(userID uint) error {
	return uc.cartRepo.ClearCart(userID)
}
