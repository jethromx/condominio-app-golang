package models

import (
	"time"

	"gorm.io/gorm"
)

// Reservation define la estructura de una reservación
type Reservation struct {
	gorm.Model
	ID              uint      `gorm:"primaryKey"`
	ResidentID      uint      `gorm:"not null"`
	ReservationDate time.Time `gorm:"not null"`
	StartTime       time.Time `gorm:"not null"`
	EndTime         time.Time `gorm:"not null"`
	Status          string    `gorm:"size:50;default:'Pending'"`
}
