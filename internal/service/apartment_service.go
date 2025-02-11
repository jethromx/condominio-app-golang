package service

import (
	"errors"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
	"github.com/gofiber/fiber/v2/log"
)

type ApartmentService interface {
	CreateApartment(apartment *models.Apartment) error
	GetApartmentByID(id uint, params map[string]interface{}) (*models.Apartment, error)
	GetAllApartments(page, pageSize int, params map[string]interface{}) ([]models.Apartment, int64, error)
	ValidateApartment(idBuilding int, idCondominium int) (*models.Condominium, error)
	UpdateApartment(apartment *models.Apartment) error
	DeleteApartment(id uint) error
}

type apartmentService struct {
	apartmentRepo repository.ApartmentRepository
	buildingRepo  repository.BuildingRepository
}

func NewApartmentService(apartmentRepo repository.ApartmentRepository, buildingRepo repository.BuildingRepository) ApartmentService {
	return &apartmentService{apartmentRepo, buildingRepo}
}

func (s *apartmentService) ValidateApartment(idBuilding int, idCondominium int) (*models.Condominium, error) {

	var condominium *models.Condominium
	var err error

	condominium, err = s.apartmentRepo.ValidationsApartment(idBuilding, idCondominium)
	if err != nil {
		return nil, err
	}

	if condominium == nil {
		log.Debug("Building or condominium does not exist")
		return nil, errors.New("building or condominium does not exist")
	} else {
		log.Debug("Building and condominium exist")
		return condominium, nil
	}

}

func (s *apartmentService) CreateApartment(apartment *models.Apartment) error {
	var err error
	var aux *models.Apartment
	var auxBuilding *models.Building

	if apartment == nil {
		return models.ErrNilApartment

	}

	if apartment.BuildingID != 0 {
		auxBuilding, err = s.buildingRepo.FindByID(auxBuilding, apartment.BuildingID)
		if err != nil {
			log.Debug("Error finding apartment by name", err)
			return errors.New("error finding building")
		}

		if auxBuilding == nil {
			return errors.New("building does not exist")
		}
	}

	aux, err = s.apartmentRepo.FindByField(aux, "name", apartment.Name)

	// Error al buscar el apartment por nombre
	if err != nil {
		log.Debug("Error finding apartment by name", err)
		return errors.New("error finding apartment")
	}

	// No existe un condominio con el mismo nombre
	if aux != nil {
		return errors.New("apartment already exists")
	}

	return s.apartmentRepo.Create(apartment)
}

func (s *apartmentService) GetApartmentByID(id uint, params map[string]interface{}) (*models.Apartment, error) {
	var apartment *models.Apartment
	var err error
	preload := params["preload"].(bool)

	if preload {
		apartment, err = s.apartmentRepo.FindByIDWithPreload(apartment, id, "Residents")
		if err != nil {
			return nil, err
		}

	} else {
		apartment, err = s.apartmentRepo.FindByID(apartment, id)
		if err != nil {
			return nil, err
		}
	}

	return apartment, nil

}

func (s *apartmentService) GetAllApartments(page, pageSize int, params map[string]interface{}) ([]models.Apartment, int64, error) {
	var apartments []models.Apartment
	var err error
	var totalRecords int64

	preload := params["preload"].(bool)
	buildingId := params["id"].(int)

	if buildingId != 0 {

		apartments, totalRecords, err = s.apartmentRepo.FindAllByBuildingID(page, pageSize, "Residents", uint(buildingId))
		if err != nil {
			return nil, 0, err
		}
		return apartments, totalRecords, nil
	}

	if preload {
		log.Debug("Preloading buildings")
		apartments, totalRecords, err = s.apartmentRepo.FindAllWithPreloadRel(page, pageSize, "Residents")
		if err != nil {
			return nil, 0, err
		}
	} else {
		log.Debug("Not preloading buildings")
		apartments, totalRecords, err = s.apartmentRepo.FindAll(page, pageSize)
		if err != nil {
			return nil, 0, err
		}
	}

	return apartments, totalRecords, nil

}

func (s *apartmentService) UpdateApartment(apartment *models.Apartment) error {
	var apartmentAux *models.Apartment
	var buildingAux *models.Building
	var err error

	buildingAux, err = s.buildingRepo.FindByID(buildingAux, apartment.BuildingID)

	if err != nil {
		return err
	}

	if buildingAux == nil {
		return errors.New("building does not exist")
	}

	apartmentAux, err = s.apartmentRepo.FindByID(apartmentAux, apartment.ID)
	if err != nil {
		return err
	}

	apartmentAux, err = s.apartmentRepo.FindByField(apartmentAux, "name", apartment.Name)
	if err != nil {
		return err
	}

	if apartmentAux != nil && apartment.ID != apartmentAux.ID {
		return errors.New("condominium already exists")
	}

	return s.apartmentRepo.Update(apartment)
}

func (s *apartmentService) DeleteApartment(id uint) error {
	var apartment *models.Apartment
	apartment, err := s.apartmentRepo.FindByID(apartment, id)
	if err != nil {
		return err
	}
	return s.apartmentRepo.Delete(apartment)
}
