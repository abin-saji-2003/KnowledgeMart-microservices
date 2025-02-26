package handlers

import (
	"context"
	"log"
	"product-service/internal/usecase"

	productProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/product-pb"
)

type ProductHandler struct {
	productUC *usecase.ProductUseCase
	productProto.UnimplementedProductServiceServer
}

func NewProductHandler(productUC *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUC: productUC}
}

func (h *ProductHandler) AddProduct(ctx context.Context, req *productProto.AddProductRequest) (*productProto.ProductResponse, error) {
	// Use correct field names from proto request
	product, err := h.productUC.AddProduct(
		req.Name,
		req.Description,
		req.ImageUrl, // ✅ Fix: Use ImageUrl instead of Image
		float64(req.Price),
		float64(req.OfferAmount),
		uint(req.SellerId),
	)
	if err != nil {
		return &productProto.ProductResponse{Status: "Failed", Message: err.Error(), Data: nil}, err
	}

	log.Print("product", product)

	// Ensure correct response mapping
	productProtoData := &productProto.Product{
		Id:           uint32(product.ID),
		Name:         product.Name,
		Description:  product.Description,
		Availability: product.Availability,
		Price:        float32(product.Price),
		OfferAmount:  float32(product.OfferAmount),
		ImageUrl:     product.Image, // ✅ Fix: Ensure ImageUrl is used
	}

	return &productProto.ProductResponse{
		Status:  "Success",
		Message: "Product successfully added to the store",
		Data:    productProtoData,
	}, nil
}

func (h *ProductHandler) GetAllProducts(ctx context.Context, req *productProto.Empty) (*productProto.ProductsResponse, error) {
	products, err := h.productUC.GetAllProducts()
	if err != nil {
		return nil, err
	}

	var productList []*productProto.Product
	for _, p := range products {
		productList = append(productList, &productProto.Product{
			Id:           uint32(p.ID),
			Name:         p.Name,
			Description:  p.Description,
			Availability: p.Availability,
			Price:        float32(p.Price),
			OfferAmount:  float32(p.OfferAmount),
			ImageUrl:     p.Image,
		})
	}

	return &productProto.ProductsResponse{
		Status:   "Success",
		Message:  "Fetched all products",
		Products: productList,
	}, nil
}

func (h *ProductHandler) GetProductById(ctx context.Context, req *productProto.GetProductRequest) (*productProto.ProductResponse, error) {
	log.Printf("Received gRPC request for product ID: %d", req.ProductId)
	product, err := h.productUC.GetProductById(uint(req.ProductId))
	if err != nil {
		return &productProto.ProductResponse{
			Status:  "error",
			Message: "Product not found",
			Data:    nil,
		}, err
	}

	return &productProto.ProductResponse{
		Status:  "success",
		Message: "Product retrieved successfully",
		Data: &productProto.Product{
			Id:           uint32(product.ID),
			Name:         product.Name,
			Description:  product.Description,
			Availability: product.Availability,
			Price:        float32(product.Price),
			OfferAmount:  float32(product.OfferAmount),
			ImageUrl:     product.Image,
		},
	}, nil
}

func (h *ProductHandler) EditProduct(ctx context.Context, req *productProto.EditProductRequest) (*productProto.ProductResponse, error) {
	// Extract optional fields
	var name, description, imageUrl *string
	var price, offerAmount *float64
	var availability *bool
	var categoryID *uint

	if req.Name != nil {
		name = req.Name
	}
	if req.Description != nil {
		description = req.Description
	}
	if req.ImageUrl != nil {
		imageUrl = req.ImageUrl
	}
	if req.Price != nil {
		priceValue := float64(req.GetPrice())
		price = &priceValue
	}
	if req.OfferAmount != nil {
		offerAmountValue := float64(req.GetOfferAmount())
		offerAmount = &offerAmountValue
	}
	if req.Availability != nil {
		availabilityValue := req.GetAvailability()
		availability = &availabilityValue
	}
	if req.CategoryId != nil {
		categoryIDValue := uint(req.GetCategoryId())
		categoryID = &categoryIDValue
	}

	product, err := h.productUC.EditProduct(
		uint(req.ProductId),
		name,
		description,
		imageUrl,
		price,
		offerAmount,
		availability,
		categoryID,
		uint(req.SellerId),
	)
	if err != nil {
		return &productProto.ProductResponse{Status: "Failed", Message: err.Error(), Data: nil}, err
	}

	// Response mapping
	productProtoData := &productProto.Product{
		Id:           uint32(product.ID),
		Name:         product.Name,
		Description:  product.Description,
		Availability: product.Availability,
		Price:        float32(product.Price),
		OfferAmount:  float32(product.OfferAmount),
		ImageUrl:     product.Image,
	}

	return &productProto.ProductResponse{
		Status:  "Success",
		Message: "Product successfully updated",
		Data:    productProtoData,
	}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *productProto.DeleteProductRequest) (*productProto.DeleteProductResponse, error) {
	log.Println("Delete request received for Product ID:", req.ProductId)

	// Call the use case to delete the product
	err := h.productUC.DeleteProduct(uint(req.ProductId), uint(req.SellerId))
	if err != nil {
		log.Println("Error deleting product:", err)
		return &productProto.DeleteProductResponse{
			Status:  "Failed",
			Message: err.Error(),
		}, err
	}

	return &productProto.DeleteProductResponse{
		Status:  "Success",
		Message: "Product successfully deleted",
	}, nil
}
