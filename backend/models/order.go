package models

import (
	"time"
)

type Order struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Currency     string    `json:"currency"`
	MerchantEmail string   `json:"merchant_email"`
	Salt         string    `json:"salt"`
	Products     []OrderProduct `json:"products" gorm:"foreignKey:OrderID"`
	TotalPrice   float64   `json:"total_price"`
	UserID       *uint     `json:"user_id"` // Nullable for guest users
	Username     string    `json:"username"` // Stores either username or "guest"
	Digest       string    `json:"digest"`
	Invoice      string    `json:"invoice" gorm:"index"`
	Status       string    `json:"status" gorm:"default:'pending'"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type OrderProduct struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Price at time of purchase
}
