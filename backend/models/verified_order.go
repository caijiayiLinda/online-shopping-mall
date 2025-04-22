package models

import "time"

type VerifiedOrder struct {
	ID           uint      `gorm:"primaryKey"`
	CreatedAt    time.Time 
	UpdatedAt    time.Time
	OrderID      uint    `gorm:"not null"`
	Invoice      string  `gorm:"uniqueIndex;not null"` 
	UserID       *uint   // Nullable for guest users
	Username     string
	Email        string  `gorm:"not null"`
	TotalPrice   float64 `gorm:"not null"`
	Currency     string  `gorm:"not null"`
	Status       string  `gorm:"not null"`
	Products     []VerifiedOrderProduct `gorm:"foreignKey:VerifiedOrderID"`
}

type VerifiedOrderProduct struct {
	ID              uint      `gorm:"primaryKey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	VerifiedOrderID uint64  `gorm:"type:bigint unsigned;not null"`
	ProductID       uint    `gorm:"not null"`
	Quantity        int     `gorm:"not null"`
	Price           float64 `gorm:"not null;type:decimal(10,2)"`
}
