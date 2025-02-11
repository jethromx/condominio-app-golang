package repository

import (
	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type CondominiumRepository interface {
	Repository[models.Condominium]
	FindAllWithPreload(page, pageSize int) ([]models.Condominium, int64, error)
}

type condominiumRepository struct {
	*BaseRepository[models.Condominium]
}

func NewCondominiumRepository(db *gorm.DB) CondominiumRepository {
	return &condominiumRepository{
		BaseRepository: NewBaseRepository[models.Condominium](db),
	}
}

func (r *BaseRepository[T]) FindAllWithPreload(page, pageSize int) ([]models.Condominium, int64, error) {
	var entities []models.Condominium
	var totalRecords int64

	offset := (page - 1) * pageSize
	if err := r.DB.Offset(offset).Limit(pageSize).Preload("Buildings").Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Model(&entities).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	return entities, totalRecords, nil
}
