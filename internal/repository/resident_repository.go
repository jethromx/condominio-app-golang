package repository

import (
	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type ResidentRepository interface {
	Repository[models.Resident]
	ValidateApartment(idBuilding int, idCondominium int, idApartment int) (*models.Apartment, error)
}

type residentRepository struct {
	*BaseRepository[models.Resident]
}

func NewResidentRepository(db *gorm.DB) ResidentRepository {
	return &residentRepository{
		BaseRepository: NewBaseRepository[models.Resident](db),
	}
}

func (r *BaseRepository[T]) ValidateApartment(idBuilding int, idCondominium int, idApartment int) (*models.Apartment, error) {
	var entity *models.Apartment

	if err := r.DB.
		Joins("JOIN buildings    b ON apartments.building_id = b.id and b.id = ?", idBuilding).
		Joins("JOIN condominiums c ON b.condominium_id = c.id and c.id = ?", idCondominium).
		Where("apartments.id = ? and apartments.deleted_at IS NULL", idApartment).Find(&entity).Error; err != nil {
		return nil, err
	}

	return entity, nil

}
