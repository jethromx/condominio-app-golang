package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID         uint    `gorm:"primaryKey"`
	Username   string  `gorm:"uniqueIndex;not null" json:"username"`
	Email      string  `gorm:"uniqueIndex;not null" json:"email"`
	Password   string  `gorm:"not null" json:"password"`
	Token      []Token `gorm:"foreignKey:UserID" json:"token"`
	Role       string  `gorm:"not null" json:"role" validate:"required,oneof=admin manager resident"` // Role del usuario
	ResidentID *uint   `gorm:"index" json:"resident_id" `                                             // Permitir valores nulos

	CreatedBy uint `gorm:"not null" json:"created_by"`
	UpdatedBy uint `gorm:"not null" json:"updated_by"`
}

// ValidateUserRole valida que el rol del usuario sea uno de los roles permitidos
func ValidateUserRole(role string) error {
	validate := validator.New()
	return validate.Var(role, "required,oneof=admin manager resident")
}
