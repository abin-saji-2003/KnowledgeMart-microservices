package models

type Cart struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint `gorm:"not null" json:"userId"`
	ProductID uint `gorm:"not null" json:"productId"`
}

type CartProduct struct {
	ProductID    uint
	Name         string
	Description  string
	Price        float64
	OfferAmount  float64
	ImageURL     string
	Availability bool
}
