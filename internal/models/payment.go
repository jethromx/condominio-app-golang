package models

import (
	"time"

	"gorm.io/gorm"
)

// Payment define la estructura de un pago
type Payment struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey"`
	ResidentID    uint      `gorm:"not null"`
	Amount        float64   `gorm:"type:decimal(10,2);not null"`
	PaymentDate   time.Time `gorm:"not null"`
	Description   string    `gorm:"size:255;not null"`
	PaymentMethod string    `gorm:"size:50;not null"`
	Status        string    `gorm:"size:50;default:'Pending'"`
}
