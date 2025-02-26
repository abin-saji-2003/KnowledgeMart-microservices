package repository

import (
	"authentication-service/internal/models"
	"gorm.io/gorm"
)

type SellerRepository interface {
	GetById(id uint) (*models.Seller, error)
	GetByUserName(name string) (*models.Seller, error)
	GetByUserId(userID uint) (*models.Seller, error)
	Create(seller *models.Seller) (*models.Seller, error)
}

type sellerRepository struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) SellerRepository {
	return &sellerRepository{db: db}
}

// GetByUserName implements SellerRepository.
func (r *sellerRepository) GetByUserName(name string) (*models.Seller, error) {
	var seller models.Seller
	err := r.db.Where("user_name = ?", name).First(&seller).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

// GetById implements SellerRepository.
func (r *sellerRepository) GetByUserId(userID uint) (*models.Seller, error) {
	var seller models.Seller
	err := r.db.Where("user_id = ?", userID).First(&seller).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

func (r *sellerRepository) Create(seller *models.Seller) (*models.Seller, error) {
	if err := r.db.Create(seller).Error; err != nil {
		return nil, err
	}
	return seller, nil
}

func (r *sellerRepository) GetById(id uint) (*models.Seller, error) {
	var seller models.Seller
	err := r.db.Where("id = ?", id).First(&seller).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}
