package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	SellerID     uint    `gorm:"not null;constraint:OnDelete:CASCADE;" json:"sellerId"`
	Name         string  `gorm:"type:varchar(255)" validate:"required" json:"name"`
	Description  string  `gorm:"type:varchar(255)" validate:"required" json:"description"`
	Availability bool    `gorm:"type:bool;default:true" json:"availability"`
	Price        float64 `gorm:"type:decimal(10,2);not null" validate:"required" json:"price"`
	OfferAmount  float64 `gorm:"type:decimal(10,2);not null" validate:"required" json:"offer_amount"`
	Image        string  `gorm:"type:varchar(255)" validate:"required" json:"image_url"`
}
