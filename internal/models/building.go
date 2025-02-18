package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNilBuilding           = errors.New("building is nil")
	ErrBuildingAlreadyExists = errors.New("building already exists")
	ErrBuildingNotFound      = errors.New("building not found")
	ErrBuildingByIdNotFound  = errors.New("building by id not found")
	ErrBuldingSameName       = errors.New("in this condominium a building already exists with the same name")
	ErrBuildingNameNotFound  = errors.New("building by name not found")
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
