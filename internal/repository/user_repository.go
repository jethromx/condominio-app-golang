package repository

import (
	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Repository[models.User]
}

type userRepository struct {
	*BaseRepository[models.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[models.User](db),
	}
}
