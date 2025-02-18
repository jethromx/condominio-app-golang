package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrResidentsNotFound = errors.New("residents not found")
)

// Resident define la estructura de un residente
type Resident struct {
	gorm.Model
	ID          uint          `gorm:"primaryKey"`
	ApartmentID uint          `gorm:"not null"`
	FirstName   string        `gorm:"size:100;not null"`
	LastName    string        `gorm:"size:100;not null"`
	Phone       string        `gorm:"size:15"`
	Email       string        `gorm:"size:100"`
	Estatus     string        `gorm:"size:50;default:'Active'"`
	Payments    []Payment     `gorm:"foreignKey:ResidentID"`
	Reservation []Reservation `gorm:"foreignKey:ResidentID"`
	UserID      uint          `gorm:"index"` // Permitir valores nulos
	//User        *User         `gorm:"foreignKey:UserID"`
}
