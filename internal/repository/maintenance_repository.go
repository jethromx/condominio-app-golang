package repository

import (
	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type MaintenanceRepository interface {
	Create(maintenance *models.Maintenance) error
	FindByID(id uint) (*models.Maintenance, error)
	FindAll() ([]models.Maintenance, error)
	Update(maintenance *models.Maintenance) error
	Delete(maintenance *models.Maintenance) error
}

type maintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) MaintenanceRepository {
	return &maintenanceRepository{db}
}

func (r *maintenanceRepository) Create(maintenance *models.Maintenance) error {
	return r.db.Create(maintenance).Error
}

func (r *maintenanceRepository) FindByID(id uint) (*models.Maintenance, error) {
	var maintenance models.Maintenance
	if err := r.db.First(&maintenance, id).Error; err != nil {
		return nil, err
	}
	return &maintenance, nil
}

func (r *maintenanceRepository) FindAll() ([]models.Maintenance, error) {
	var maintenances []models.Maintenance
	if err := r.db.Find(&maintenances).Error; err != nil {
		return nil, err
	}
	return maintenances, nil
}

func (r *maintenanceRepository) Update(maintenance *models.Maintenance) error {
	return r.db.Save(maintenance).Error
}

func (r *maintenanceRepository) Delete(maintenance *models.Maintenance) error {
	return r.db.Delete(maintenance).Error
}
