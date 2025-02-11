package models

import (
	"gorm.io/gorm"
)

// User define la estructura de un usuario
type Building struct {
	gorm.Model
	ID            uint        `gorm:"primaryKey"`
	CondominiumID uint        `gorm:"not null"`
	Name          string      `gorm:"size:100;not null"`
	Floors        int         `gorm:"not null"`
	Apartments    []Apartment `gorm:"foreignKey:BuildingID"`
	CreatedBy     string      `gorm:"size:64"`
	UpdatedBy     string      `gorm:"size:64"`
}
