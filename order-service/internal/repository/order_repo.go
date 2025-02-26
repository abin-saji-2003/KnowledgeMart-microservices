package repository

import (
	"errors"

	"order-service/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) (uint, error)
	CreateOrderItems(items models.OrderItem) error
	ClearCart(userID uint) error
	GetCartItems(userID uint) ([]models.Cart, error)
	GetOrdersByUserID(userID uint) ([]models.Order, error)
	GetOrderItemsByOrderID(orderID uint) ([]models.OrderItem, error)
	UpdateOrderSellerID(orderID, sellerID uint) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *models.Order) (uint, error) {
	if err := r.db.Create(order).Error; err != nil {
		return 0, err
	}
	return order.OrderID, nil // ✅ Return the generated OrderID
}

func (r *orderRepository) CreateOrderItems(items models.OrderItem) error {
	return r.db.Create(&items).Error
}

func (r *orderRepository) ClearCart(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Cart{}).Error
}

func (r *orderRepository) GetCartItems(userID uint) ([]models.Cart, error) {
	var cartItems []models.Cart
	err := r.db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}
	return cartItems, nil
}

func (r *orderRepository) UpdateOrderSellerID(orderID, sellerID uint) error {
	return r.db.Model(&models.Order{}).Where("order_id = ?", orderID).Update("seller_id", sellerID).Error
}

func (r *orderRepository) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetOrderItemsByOrderID(orderID uint) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.db.Where("order_id = ?", orderID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
