package repository

import (
	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type ResidentRepository interface {
	Repository[models.Resident]
}

type residentRepository struct {
	*BaseRepository[models.Resident]
}

func NewResidentRepository(db *gorm.DB) ResidentRepository {
	return &residentRepository{
		BaseRepository: NewBaseRepository[models.Resident](db),
	}
}
