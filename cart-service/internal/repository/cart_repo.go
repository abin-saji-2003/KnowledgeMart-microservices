package repository

import (
	"cart-service/internal/models"
	"errors"

	"gorm.io/gorm"
)

type CartRepository interface {
	AddtoCart(cart *models.Cart) error
	RemoveFromCart(productId, userId uint) error
	GetProductsFromCart(id uint) ([]models.Cart, error)
	ClearCart(userID uint) error
	IsProductInCart(userId, productId uint) (bool, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) ClearCart(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Cart{}).Error
}

// AddtoCart implements CartRepository.
func (c *cartRepository) AddtoCart(cart *models.Cart) error {
	return c.db.Create(cart).Error
}

// GetProductsFromCart implements CartRepository.
func (c *cartRepository) GetProductsFromCart(userID uint) ([]models.Cart, error) {
	var carts []models.Cart
	if err := c.db.Where("user_id = ?", userID).Find(&carts).Error; err != nil {
		return nil, err
	}
	return carts, nil
}

// RemoveFromCart implements CartRepository.
func (c *cartRepository) RemoveFromCart(productId, userId uint) error {
	// First, check if the product exists in the cart
	var cartItem models.Cart
	if err := c.db.Where("product_id = ? AND user_id = ?", productId, userId).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found in cart")
		}
		return err
	}

	// If found, delete it
	if err := c.db.Delete(&cartItem).Error; err != nil {
		return errors.New("failed to remove product from cart")
	}

	return nil
}

func (c *cartRepository) IsProductInCart(userId, productId uint) (bool, error) {
	var count int64
	err := c.db.Model(&models.Cart{}).
		Where("user_id = ? AND product_id = ?", userId, productId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
