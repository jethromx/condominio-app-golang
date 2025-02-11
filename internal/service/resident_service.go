package service

import (
	"errors"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
	"github.com/gofiber/fiber/v2/log"
)

type ResidentService interface {
	CreateResident(apartment *models.Resident) error
	GetResidentByID(id uint) (*models.Resident, error)
	GetAllResidents(page, pageSize int) ([]models.Resident, int64, error)
	UpdateResident(iapartment *models.Resident) error
	DeleteResident(id uint) error
	ValidateResident(id uint) error
}

type residentService struct {
	residentRepo repository.ResidentRepository
}

func NewResidentService(residentRepo repository.ResidentRepository) ResidentService {
	return &residentService{residentRepo}
}

func (s *residentService) CreateResident(resident *models.Resident) error {
	if resident == nil {
		log.Debug("resident is nil")
		return errors.New("error input data")

	}

	// TODO Validar que el residente tenga un departamento asignado

	/*
		aux, err := s.residentRepo.FindByField(resident, "first_name", resident.FirstName)
		if err != nil {
			log.Debug("Error finding resident by name", err)
			return errors.New("error finding resident by name")
		}

		// No existe un condominio con el mismo nombre
		if aux != nil {
			return errors.New("record already exists")
		} */

	return s.residentRepo.Create(resident)
}

func (s *residentService) GetResidentByID(id uint) (*models.Resident, error) {
	var resident *models.Resident
	return s.residentRepo.FindByID(resident, id)
}

func (s *residentService) GetAllResidents(page, pageSize int) ([]models.Resident, int64, error) {
	return s.residentRepo.FindAll(page, pageSize)
}

func (s *residentService) UpdateResident(resident *models.Resident) error {
	var err error
	var entityAux *models.Resident

	err = s.ValidateResident(resident.ID)
	if err != nil {
		return err
	}

	//TODO validar que exista el usuario

	entityAux, err = s.residentRepo.FindByID(entityAux, resident.ID)
	if err != nil {
		return err
	}

	/*
		entityAux, err := s.residentRepo.FindByField(resident, "first_name", resident.FirstName)

		if err != nil {
			return err
		}*/

	if entityAux != nil && resident.ID != entityAux.ID {
		return errors.New("record already exists")
	}

	return s.residentRepo.Update(resident)
}

func (s *residentService) DeleteResident(id uint) error {
	var resident *models.Resident
	resident, err := s.residentRepo.FindByID(resident, id)
	if err != nil {
		return err
	}
	return s.residentRepo.Delete(resident)
}

func (s *residentService) ValidateResident(id uint) error {
	var entity *models.Resident
	entity, err := s.residentRepo.FindByID(entity, id)

	if err != nil {
		log.Debug("Error finding record by ID", err)
		return errors.New("error finding record by ID")
	}

	if entity == nil {
		return errors.New("record does not exist")
	}

	return nil

}
