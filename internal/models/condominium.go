package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNilCondominium           = errors.New("condominium is nil")
	ErrCondominiumAlreadyExists = errors.New("condominium already exists")
	ErrCondominiumNotFound      = errors.New("condominium not found")
	ErrCondominimByIdNotFound   = errors.New("condominium by id not found")
)

// Condominium define la estructura de un condominio
type Condominium struct {
	gorm.Model
	ID        uint       `gorm:"primaryKey"`
	Name      string     `gorm:"size:100;not null"`
	Address   string     `gorm:"size:255;not null"`
	Phone     string     `gorm:"size:15"`
	Email     string     `gorm:"size:100"`
	ZIPCode   string     `gorm:"size:10"`
	Buildings []Building `gorm:"foreignKey:CondominiumID"`
	//Administradores []User `gorm:"many2many:condominio_admins;"`

	CreatedBy string `gorm:"size:64"`
	UpdatedBy string `gorm:"size:64"`
}

// TableName especifica el nombre de la tabla para el modelo Condominium
func (Condominium) TableName() string {
	return "condominiums"
}
