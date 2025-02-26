package models

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"type:varchar(255);unique" validate:"required,email" json:"email"`
	Password string `gorm:"type:varchar(255)" validate:"required" json:"password"`
}

type User struct {
	gorm.Model
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string  `gorm:"type:varchar(255)" validate:"required" json:"name"`
	Email        string  `gorm:"type:varchar(255);unique" validate:"email" json:"email"`
	PhoneNumber  string  `gorm:"type:varchar(255);unique" validate:"number" json:"phone_number"`
	Picture      string  `gorm:"type:text" json:"picture"`
	WalletAmount float64 `gorm:"column:wallet_amount;type:double precision" json:"wallet_amount"`
	Password     string  `gorm:"type:varchar(255)" validate:"required" json:"password"`
	Blocked      bool    `gorm:"type:bool" json:"blocked"`
	// OTP          uint64
	// OTPExpiry    time.Time
	// IsVerified   bool   `gorm:"type:bool" json:"verified"`
	LoginMethod string `gorm:"type:varchar(50)" json:"login_method"`
}

type Seller struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint   `gorm:"not null;constraint:OnDelete:CASCADE;" json:"userId"`
	UserName    string `gorm:"type:varchar(255)" validate:"required" json:"name"`
	Password    string `gorm:"type:varchar(255)" validate:"required" json:"password"`
	Description string `gorm:"type:varchar(255)" validate:"required" json:"description"`
}
