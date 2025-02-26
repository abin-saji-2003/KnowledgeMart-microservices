package usecase

import (
	"errors"
	"product-service/internal/models"
	"product-service/internal/repository"
)

type ProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

func (pu *ProductUseCase) AddProduct(name, description, image string, price float64, offerAmount float64, sellerId uint) (*models.Product, error) {
	if offerAmount > price {
		return nil, errors.New("offer amount should be less than actual price")
	}

	if price <= 0 {
		return nil, errors.New("product price must be a positive integer")
	}

	newProduct := &models.Product{
		Name:         name,
		Description:  description,
		Availability: true,
		Price:        price,
		OfferAmount:  offerAmount,
		Image:        image,
		SellerID:     sellerId,
	}

	if err := pu.productRepo.AddProduct(newProduct); err != nil {
		return nil, errors.New("failed to create product")
	}

	return newProduct, nil
}

func (pu *ProductUseCase) GetAllProducts() ([]models.Product, error) {
	return pu.productRepo.GetAllProducts()
}

func (pu *ProductUseCase) GetProductById(id uint) (*models.Product, error) {
	return pu.productRepo.GetProductById(id)
}

func (pu *ProductUseCase) EditProduct(productID uint, name, description, image *string, price, offerAmount *float64, availability *bool, categoryID *uint, sellerID uint) (*models.Product, error) {
	// Fetch existing product
	product, err := pu.productRepo.GetProductById(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Ensure only the seller who created it can edit
	if product.SellerID != sellerID {
		return nil, errors.New("unauthorized to edit this product")
	}

	// Update fields if provided
	if name != nil {
		product.Name = *name
	}
	if description != nil {
		product.Description = *description
	}
	if image != nil {
		product.Image = *image
	}
	if price != nil {
		if *price <= 0 {
			return nil, errors.New("product price must be a positive integer")
		}
		product.Price = *price
	}
	if offerAmount != nil {
		if *offerAmount > product.Price {
			return nil, errors.New("offer amount should be less than actual price")
		}
		product.OfferAmount = *offerAmount
	}
	if availability != nil {
		product.Availability = *availability
	}

	// Save changes
	if err := pu.productRepo.EditProduct(product); err != nil {
		return nil, errors.New("failed to update product")
	}

	return product, nil
}

func (pu *ProductUseCase) DeleteProduct(productID, sellerID uint) error {
	// Check if the product exists and belongs to the seller
	product, err := pu.productRepo.GetProductById(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Ensure that only the seller who owns the product can delete it
	if product.SellerID != sellerID {
		return errors.New("unauthorized: you can only delete your own products")
	}

	// Perform deletion
	err = pu.productRepo.DeleteProduct(productID)
	if err != nil {
		return errors.New("failed to delete product")
	}

	return nil
}
