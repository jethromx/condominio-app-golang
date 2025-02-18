package repository

import (
	"com.mx/crud/internal/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ApartmentRepository interface {
	Repository[models.Apartment]
	FindAllByBuildingID(page, pageSize int, preload string, id uint) ([]models.Apartment, int64, error)
	ValidationsApartment(idBuilding int, idCondominium int) (*models.Condominium, error)
	FindAllByCondominiumID(page, pageSize int, preload string, id uint) ([]models.Apartment, int64, error)
	ValidateResidents(ids []int) error
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
	log.Debug("ID Building: ", id)
	offset := (page - 1) * pageSize
	if err := r.DB.Offset(offset).Limit(pageSize).Preload(preload).Where(&models.Apartment{BuildingID: id}).Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	if err := r.DB.Model(&models.Apartment{}).Where(&models.Apartment{BuildingID: id}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	return entities, totalRecords, nil
}

func (r *BaseRepository[T]) FindAllByCondominiumID(page, pageSize int, preload string, id uint) ([]models.Apartment, int64, error) {
	var entities []models.Apartment
	var totalRecords int64
	offset := (page - 1) * pageSize
	if err := r.DB.
		Offset(offset).Limit(pageSize).
		Joins("JOIN buildings b ON apartments.building_id = b.id ").
		Joins("JOIN condominiums c ON b.condominium_id = c.id and c.id= ?", id).Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	if err := r.DB.
		Joins("JOIN buildings b ON apartments.building_id = b.id ").
		Joins("JOIN condominiums c ON b.condominium_id = c.id and c.id= ?", id).Find(&entities).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}
	log.Debug("Total Records: ", totalRecords)

	return entities, totalRecords, nil
}

func (r *BaseRepository[T]) ValidationsApartment(idBuilding int, idCondominium int) (*models.Condominium, error) {
	var condominium *models.Condominium

	if err := r.DB.Joins("JOIN buildings b ON b.condominium_id = condominiums.id and b.id = ?", idBuilding).Where("condominiums.id = ? and condominiums.deleted_at IS NULL", idCondominium).Find(&condominium).Error; err != nil {
		return nil, err
	}

	return condominium, nil

}

func (r *BaseRepository[T]) ValidateResidents(ids []int) error {
	var count int64

	log.Debug(ids)
	if err := r.DB.Model(&models.Resident{}).Where("id IN (?)", ids).Count(&count).Error; err != nil {
		return err
	}

	log.Debug(count)
	if count != int64(len(ids)) {
		return models.ErrResidentsNotFound
	}

	return nil
}
