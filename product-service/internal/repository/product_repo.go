package repository

import (
	"gorm.io/gorm"
	"product-service/internal/models"
)

type ProductRepository interface {
	AddProduct(product *models.Product) error
	EditProduct(product *models.Product) error
	DeleteProduct(id uint) error
	GetProductById(id uint) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// AddProduct inserts a new product into the database.
func (r *productRepository) AddProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

// EditProduct updates an existing product.
func (r *productRepository) EditProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

// DeleteProduct removes a product from the database.
func (r *productRepository) DeleteProduct(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// GetProductById retrieves a product by its ID.
func (r *productRepository) GetProductById(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAllProducts retrieves all products from the database.
func (r *productRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
