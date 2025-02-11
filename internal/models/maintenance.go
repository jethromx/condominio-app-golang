package models

import (
	"time"

	"gorm.io/gorm"
)

// Maintenance define la estructura de un mantenimiento
type Maintenance struct {
	gorm.Model
	ID             uint      `gorm:"primaryKey"`
	ApartmentID    uint      `gorm:"not null"`
	Description    string    `gorm:"type:text;not null"`
	RequestDate    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status         string    `gorm:"size:50;default:'Pending'"`
	ResolutionDate time.Time
}
