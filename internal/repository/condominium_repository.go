package repository

import (
	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type CondominiumRepository interface {
	Repository[models.Condominium]
}

type condominiumRepository struct {
	*BaseRepository[models.Condominium]
}

func NewCondominiumRepository(db *gorm.DB) CondominiumRepository {
	return &condominiumRepository{
		BaseRepository: NewBaseRepository[models.Condominium](db),
	}
}
