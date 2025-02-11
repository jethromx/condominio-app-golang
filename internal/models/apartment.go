package models

import (
	"errors"

	"gorm.io/gorm"
)

// Apartment define la estructura de un departamento
type Apartment struct {
	gorm.Model
	ID                  uint          `gorm:"primaryKey"`
	BuildingID          uint          `gorm:"not null"`
	Name                string        `gorm:"size:100;not null"`
	Number              string        `gorm:"size:10;not null"`
	Floor               int           `gorm:"not null"`
	Status              string        `gorm:"size:50;default:'Available'"`
	Residents           []Resident    `gorm:"foreignKey:ApartmentID"`
	MaintenanceRequests []Maintenance `gorm:"foreignKey:ApartmentID"`
}

var (
	ErrNilApartment = errors.New("apartment is nil")
)
