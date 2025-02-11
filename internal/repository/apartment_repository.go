package repository

import (
	"github.com/gofiber/fiber/v2/log"

	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type ApartmentRepository interface {
	Repository[models.Apartment]
	FindAllByBuildingID(page, pageSize int, preload string, id uint) ([]models.Apartment, int64, error)
	ValidationsApartment(idBuilding int, idCondominium int) (*models.Condominium, error)
}

type apartmentRepository struct {
	*BaseRepository[models.Apartment]
}

func NewApartmentRepository(db *gorm.DB) ApartmentRepository {
	return &apartmentRepository{
		BaseRepository: NewBaseRepository[models.Apartment](db),
	}
}

func (r *BaseRepository[T]) FindAllByBuildingID(page, pageSize int, preload string, id uint) ([]models.Apartment, int64, error) {
	var entities []models.Apartment
	var totalRecords int64

	offset := (page - 1) * pageSize
	if err := r.DB.Offset(offset).Limit(pageSize).Preload(preload).Where(&models.Apartment{BuildingID: id}).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Model(&entities).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	return entities, totalRecords, nil
}

func (r *BaseRepository[T]) ValidationsApartment(idBuilding int, idCondominium int) (*models.Condominium, error) {
	var condominium *models.Condominium

	if err := r.DB.Joins("JOIN buildings b ON b.condominium_id = condominiums.id and b.id = ?", idBuilding).Where("condominiums.id = ?", idCondominium).Find(&condominium).Error; err != nil {
		log.Debug("Error ", err)
		return nil, err
	}

	return condominium, nil

}
