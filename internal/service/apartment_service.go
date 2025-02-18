package service

import (
	"errors"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2/log"
)

type ApartmentService interface {
	CreateApartment(apartment *models.Apartment) error
	GetApartmentByID(id uint, params map[string]interface{}) (*models.Apartment, error)
	GetAllApartments(page, pageSize int, params map[string]interface{}) ([]models.Apartment, int64, error)
	GetAllApartmentsByCondominium(page, pageSize int, condominiumID uint) ([]models.Apartment, int64, error)

	ValidateApartment(idBuilding int, idCondominium int) (*models.Condominium, error)
	ValidateResidents(ids []int) error
	UpdateApartment(apartment *models.Apartment) error
	DeleteApartment(id uint) error
}

type apartmentService struct {
	apartmentRepo repository.ApartmentRepository
	buildingRepo  repository.BuildingRepository
}

var (
	RESIDENTS = "Residents"
)

func NewApartmentService(apartmentRepo repository.ApartmentRepository, buildingRepo repository.BuildingRepository) ApartmentService {
	return &apartmentService{apartmentRepo, buildingRepo}
}

func (s *apartmentService) ValidateApartment(idBuilding int, idCondominium int) (*models.Condominium, error) {

	var condominium *models.Condominium
	var err error

	// Validar que el condominio exista y esta relacionado con el edificio
	condominium, err = s.apartmentRepo.ValidationsApartment(idBuilding, idCondominium)

	return condominium, err

}

func (s *apartmentService) GetAllApartmentsByCondominium(page, pageSize int, condominiumID uint) ([]models.Apartment, int64, error) {

	return s.apartmentRepo.FindAllByCondominiumID(page, pageSize, "", condominiumID)
}

func (s *apartmentService) ValidateResidents(ids []int) error {

	var err error

	// Validar que existan los residentes
	err = s.apartmentRepo.ValidateResidents(ids)

	return err

}

func (s *apartmentService) CreateApartment(apartment *models.Apartment) error {
	var err error
	var aux = &models.Apartment{}
	var auxBuilding = &models.Building{}

	if apartment == nil {
		return models.ErrNilApartment

	}

	if apartment.BuildingID != 0 {
		if err = s.buildingRepo.FindID(auxBuilding, apartment.BuildingID); err != nil {
			return errors.New(utils.MsgBuildingNotFound)
		}

		if auxBuilding.ID == 0 {
			return errors.New(utils.MsgBuildingNotExist)
		}
	}

	if err = s.apartmentRepo.FindField(aux, NAME, apartment.Name); err != nil {
		return errors.New(utils.MsgApartmentNotFound)
	}

	// No existe un condominio con el mismo nombre
	if aux.ID != 0 {
		return errors.New(utils.MsgApartmentAlreadyExists)
	}

	return s.apartmentRepo.Create(apartment)
}

func (s *apartmentService) GetApartmentByID(id uint, params map[string]interface{}) (*models.Apartment, error) {
	var apartment = &models.Apartment{}
	var err error
	preload := params[PRELOAD].(bool)

	if preload {
		err = s.apartmentRepo.FindByIDPreload(apartment, id, RESIDENTS)
		if err != nil {
			return nil, err
		}

	} else {
		err = s.apartmentRepo.FindID(apartment, id)
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

	log.Debug(buildingId)
	log.Debug(preload)

	if buildingId == 0 {

		if preload {
			log.Debug("loading buildings")
			apartments, totalRecords, err = s.apartmentRepo.FindAllWithPreloadRel(page, pageSize, "Residents")
			if err != nil {
				return nil, 0, err
			}
		} else {
			log.Debug("loading buildings without preloading")
			apartments, totalRecords, err = s.apartmentRepo.FindAllWithPreloadRel(page, pageSize, "")
			if err != nil {
				return nil, 0, err
			}
		}
		//return apartments, totalRecords, nil
	} else {
		if preload {
			log.Debug("Preloading buildings")
			apartments, totalRecords, err = s.apartmentRepo.FindAllByBuildingID(page, pageSize, "Residents", uint(buildingId))
			if err != nil {
				return nil, 0, err
			}
		} else {
			log.Debug("Not preloading buildings")
			apartments, totalRecords, err = s.apartmentRepo.FindAllByBuildingID(page, pageSize, "", uint(buildingId))
			if err != nil {
				return nil, 0, err
			}
		}

	}

	return apartments, totalRecords, nil

}

func (s *apartmentService) UpdateApartment(apartment *models.Apartment) error {
	var apartmentAux = &models.Apartment{}
	var buildingAux = &models.Building{}

	var err error

	err = s.buildingRepo.FindID(buildingAux, apartment.BuildingID)

	if err != nil {
		return err
	}

	if buildingAux.ID == 0 {
		return errors.New("building does not exist")
	}

	err = s.apartmentRepo.FindID(apartmentAux, apartment.ID)
	if err != nil {
		return err
	}

	err = s.apartmentRepo.FindField(apartmentAux, NAME, apartment.Name)
	if err != nil {
		return err
	}

	if apartment.ID != apartmentAux.ID {
		return errors.New("condominium already exists")
	}

	return s.apartmentRepo.UpdatePreLoad(apartment, "Residents")
}

func (s *apartmentService) DeleteApartment(id uint) error {
	var apartment = &models.Apartment{}
	err := s.apartmentRepo.FindID(apartment, id)
	if err != nil {
		return err
	}
	return s.apartmentRepo.Delete(apartment)
}
