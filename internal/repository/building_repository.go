package repository

import (
	"errors"

	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type BuildingRepository interface {
	Repository[models.Building]
	FindBuildings(entities []models.Building, ID int64) ([]models.Building, error)
	FindBuildingsByCondominium(page, pageSize int, ID int) ([]models.Building, int64, error)
	FindBuildingsByCondominiumPreload(page, pageSize int, ID int) ([]models.Building, int64, error)
	ValidateBuilding(idCondominuim int, field string) (*models.Building, error)
	ValidateBuildingByIDs(idCondominuim int, idBuilding int) (*models.Building, error)
}

type buildingRepository struct {
	*BaseRepository[models.Building]
}

func NewBuildingRepository(db *gorm.DB) BuildingRepository {
	return &buildingRepository{
		BaseRepository: NewBaseRepository[models.Building](db),
	}
}

func (r *BaseRepository[T]) FindBuildings(entities []models.Building, ID int64) ([]models.Building, error) {

	if err := r.DB.Find(&entities, &models.Building{CondominiumID: uint(ID)}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return entities, nil
}

func (r *BaseRepository[T]) FindBuildingsByCondominium(page, pageSize int, ID int) ([]models.Building, int64, error) {
	var entities []models.Building
	var totalRecords int64

	offset := (page - 1) * pageSize
	if err := r.DB.Offset(offset).Limit(pageSize).Find(&entities, &models.Building{CondominiumID: uint(ID)}).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Model(&entities).Where(&models.Building{CondominiumID: uint(ID)}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	return entities, totalRecords, nil
}

func (r *BaseRepository[T]) FindBuildingsByCondominiumPreload(page, pageSize int, ID int) ([]models.Building, int64, error) {
	var entities []models.Building
	var totalRecords int64

	offset := (page - 1) * pageSize
	if err := r.DB.Offset(offset).Limit(pageSize).Preload("Apartments").Find(&entities, &models.Building{CondominiumID: uint(ID)}).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Model(&entities).Where(&models.Building{CondominiumID: uint(ID)}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	return entities, totalRecords, nil
}

func (r *BaseRepository[T]) ValidateBuilding(idCondominuim int, field string) (*models.Building, error) {
	var entity models.Building
	if err := r.DB.Where("condominium_id = ? AND name = ?", idCondominuim, field).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.New("error finding building by name")
	}

	return &entity, nil
}

func (r *BaseRepository[T]) ValidateBuildingByIDs(idCondominuim int, idBuilding int) (*models.Building, error) {
	var entity models.Building
	var condominium models.Condominium
	condominium.ID = uint(idCondominuim)
	entity.ID = uint(idBuilding)
	entity.CondominiumID = uint(idCondominuim)

	if err := r.DB.First(&condominium).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.New("error finding condominium building ")
	}

	if err := r.DB.First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.New("error finding building")
	}

	return &entity, nil
}
