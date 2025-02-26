package models

import "time"

type Order struct {
	OrderID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                 uint      `gorm:"not null" json:"user_id"`
	CouponCode             string    `json:"coupon_code"`
	CouponDiscountAmount   float64   `validate:"required,number" json:"coupon_discount_amount"`
	CategoryDiscountAmount float64   `validate:"required,number" json:"category_discount_amount"`
	ProductOfferAmount     float64   `validate:"required,number" json:"product_offer_amount"`
	DeliveryCharge         float64   `validate:"number" json:"delivery_charge"`
	TotalAmount            float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	FinalAmount            float64   `validate:"required,number" json:"final_amount"`
	PaymentMethod          string    `gorm:"type:varchar(100)" validate:"required" json:"payment_method"`
	PaymentStatus          string    `gorm:"type:varchar(100)" validate:"required" json:"payment_status"`
	OrderedAt              time.Time `gorm:"autoCreateTime" json:"ordered_at"`
	// ShippingAddress        ShippingAddress `gorm:"embedded" json:"shipping_address"`
	SellerID           uint   `gorm:"not null" json:"seller_id"`
	Status             string `gorm:"type:varchar(100);default:'pending'" json:"status"`
	FailedPaymentCount int    `gorm:"default:0" json:"failed_Payment_count"`
}

type OrderItem struct {
	OrderItemID         uint    `gorm:"primaryKey;autoIncrement" json:"orderItemId"`
	OrderID             uint    `gorm:"not null" json:"orderId"`
	UserID              uint    `gorm:"not null" json:"userId"`
	ProductID           uint    `gorm:"not null" json:"productId"`
	SellerID            uint    `gorm:"not null" json:"sellerId"`
	Price               float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	ProductOfferAmount  float64 `json:"product_offer_amount"`
	CategoryOfferAmount float64 `json:"category_offer_amount"`
	OtherOffers         float64 `json:"other_offers"`
	FinalAmount         float64 `json:"final_amount"`
	Status              string  `gorm:"type:varchar(100);default:'pending'" json:"status"`
}

type Cart struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint    `gorm:"not null" json:"user_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`
}

type OrderResponse struct {
	OrderID        uint                `json:"id"`
	UserID         uint                `json:"user_id"`
	DeliveryCharge float64             `json:"delivery_charge"`
	TotalAmount    float64             `json:"total_amount"`
	FinalAmount    float64             `json:"final_amount"`
	PaymentMethod  string              `json:"payment_method"`
	PaymentStatus  string              `json:"payment_status"`
	OrderedAt      time.Time           `json:"ordered_at"`
	SellerID       uint                `json:"seller_id"`
	Status         string              `json:"status"`
	OrderItems     []OrderItemResponse `json:"order_items"`
}

type OrderItemResponse struct {
	OrderID     uint    `gorm:"not null" json:"orderId"`
	UserID      uint    `gorm:"not null" json:"userId"`
	ProductID   uint    `gorm:"not null" json:"productId"`
	SellerID    uint    `gorm:"not null" json:"sellerId"`
	Name        string  `gorm:"type:varchar(255)" validate:"required" json:"name"`
	Description string  `gorm:"type:varchar(255)" validate:"required" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	FinalAmount float64 `json:"final_amount"`
	Status      string  `gorm:"type:varchar(100);default:'pending'" json:"status"`
	Image       string  `gorm:"type:varchar(255)" validate:"required" json:"image_url"`
}

// type ShippingAddress struct {
// 	StreetName   string `gorm:"type:varchar(255)" json:"street_name"`
// 	StreetNumber string `gorm:"type:varchar(255)" json:"street_number"`
// 	City         string `gorm:"type:varchar(255)" json:"city"`
// 	State        string `gorm:"type:varchar(255)" json:"state"`
// 	PinCode      string `gorm:"type:varchar(20)" json:"pincode"`
// 	PhoneNumber  string `gorm:"type:varchar(20)" json:"phonenumber"`
// }
