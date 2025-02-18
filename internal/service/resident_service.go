package service

import (
	"errors"
	"strconv"

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
	ValidateApartment(idBuilding int, idCondominium int, idApartment int, idUser int) error
	ValidateApartmentSU(idBuilding int, idCondominium int, idApartment int) error
}

type residentService struct {
	residentRepo repository.ResidentRepository
	userRepo     repository.UserRepository
}

func NewResidentService(residentRepo repository.ResidentRepository, userRepo repository.UserRepository) ResidentService {
	return &residentService{residentRepo, userRepo}
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
	var resident = &models.Resident{}
	err := s.residentRepo.FindID(resident, id)
	if err != nil {
		return nil, err
	}
	return resident, nil
}

func (s *residentService) GetAllResidents(page, pageSize int) ([]models.Resident, int64, error) {
	return s.residentRepo.FindAll(page, pageSize)
}

func (s *residentService) UpdateResident(resident *models.Resident) error {
	var err error
	var entityAux = &models.Resident{}

	err = s.residentRepo.FindID(entityAux, resident.ID)
	if err != nil {
		return err
	}

	if entityAux.ID != 0 && resident.ID != entityAux.ID {
		return errors.New("record already exists")
	}

	return s.residentRepo.Update(resident)
}

func (s *residentService) DeleteResident(id uint) error {
	var resident = &models.Resident{}
	err := s.residentRepo.FindID(resident, id)
	if err != nil {
		return err
	}
	return s.residentRepo.Delete(resident)
}

func (s *residentService) ValidateResident(id uint) error {
	var entity *models.Resident
	err := s.residentRepo.FindID(entity, id)

	if err != nil {
		log.Debug("Error finding record by ID", err)
		return errors.New("error finding record by ID")
	}

	if entity == nil {
		return errors.New("record does not exist")
	}

	return nil

}

func (s *residentService) ValidateApartment(idBuilding int, idCondominium int, idApartment int, idUser int) error {
	var apartment *models.Apartment
	var user = models.User{}
	var resident = &models.Resident{}
	apartment, err := s.residentRepo.ValidateApartment(idBuilding, idCondominium, idApartment)

	if err != nil {
		log.Debug("Error validating apartment", err)
		return errors.New("error validating apartment")
	}

	if apartment == nil {
		return errors.New("apartment does not exist")
	}

	if err = s.userRepo.FindID(&user, uint(idUser)); err != nil {
		log.Debug("Error finding user", err)
		return errors.New("error finding user")
	}

	if user.ID == 0 {
		return errors.New("user does not exist")
	}

	if err = s.residentRepo.FindField(resident, "user_id", strconv.Itoa(idUser)); err != nil {
		log.Debug("Error finding resident", err)
		return errors.New("error finding resident")
	}

	log.Debug("Resident", resident.ID)
	if resident.ID != 0 {
		return errors.New("user already exists asociated to a resident")
	}

	return nil

}

func (s *residentService) ValidateApartmentSU(idBuilding int, idCondominium int, idApartment int) error {
	var apartment *models.Apartment
	var resident = &models.Resident{}
	apartment, err := s.residentRepo.ValidateApartment(idBuilding, idCondominium, idApartment)

	if err != nil {
		log.Debug("Error validating apartment", err)
		return errors.New("error validating apartment")
	}

	if apartment == nil {
		return errors.New("apartment does not exist")
	}

	if resident.ID != 0 {
		return errors.New("user already exists asociated to a resident")
	}

	return nil

}
