package service

import (
	"time"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
)

type MaintenanceService interface {
	CreateMaintenance(description string, date time.Time, cost float64, condominiumID uint) (*models.Maintenance, error)
	GetMaintenanceByID(id uint) (*models.Maintenance, error)
	GetAllMaintenances() ([]models.Maintenance, error)
	UpdateMaintenance(id uint, description string, date time.Time, cost float64, condominiumID uint) (*models.Maintenance, error)
	DeleteMaintenance(id uint) error
}

type maintenanceService struct {
	maintenanceRepo repository.MaintenanceRepository
}

func NewMaintenanceService(maintenanceRepo repository.MaintenanceRepository) MaintenanceService {
	return &maintenanceService{maintenanceRepo}
}

func (s *maintenanceService) CreateMaintenance(description string, date time.Time, cost float64, condominiumID uint) (*models.Maintenance, error) {
	maintenance := &models.Maintenance{
		Description: description,
		//Date:          date,
		//Cost:          cost,
		//CondominiumID: condominiumID,
	}
	if err := s.maintenanceRepo.Create(maintenance); err != nil {
		return nil, err
	}
	return maintenance, nil
}

func (s *maintenanceService) GetMaintenanceByID(id uint) (*models.Maintenance, error) {
	return s.maintenanceRepo.FindByID(id)
}

func (s *maintenanceService) GetAllMaintenances() ([]models.Maintenance, error) {
	return s.maintenanceRepo.FindAll()
}

func (s *maintenanceService) UpdateMaintenance(id uint, description string, date time.Time, cost float64, condominiumID uint) (*models.Maintenance, error) {
	maintenance, err := s.maintenanceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	maintenance.Description = description
	//maintenance.Date = date
	//maintenance.Cost = cost
	//maintenance.CondominiumID = condominiumID
	if err := s.maintenanceRepo.Update(maintenance); err != nil {
		return nil, err
	}
	return maintenance, nil
}

func (s *maintenanceService) DeleteMaintenance(id uint) error {
	maintenance, err := s.maintenanceRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.maintenanceRepo.Delete(maintenance)
}
