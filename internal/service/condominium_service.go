package service

import (
	"errors"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
	"github.com/gofiber/fiber/v2/log"
)

type CondominiumService interface {
	CreateCondominium(condominium *models.Condominium) error
	GetCondominiumByID(id uint, params map[string]interface{}) (*models.Condominium, error)
	GetAllCondominiums(page, pageSize int, params map[string]interface{}) ([]models.Condominium, int64, error)
	UpdateCondominium(condominium *models.Condominium) error
	DeleteCondominium(id uint) error
}

type condominiumService struct {
	condominiumRepo repository.CondominiumRepository
	buildingRepo    repository.BuildingRepository
}

func NewCondominiumService(condominiumRepo repository.CondominiumRepository, buildingRepo repository.BuildingRepository) CondominiumService {
	return &condominiumService{condominiumRepo, buildingRepo}
}

func (s *condominiumService) CreateCondominium(condominium *models.Condominium) error {

	if condominium == nil {

		return models.ErrNilCondominium
	}

	aux, err := s.condominiumRepo.FindByField(condominium, "name", condominium.Name)

	// Existe un condominio con el mismo nombre
	if err != nil {
		log.Debug("Error finding condominium by name", err)
		return errors.New("error input data")
	}

	// No existe un condominio con el mismo nombre
	if aux != nil {
		return errors.New("condominium already exists")
	}

	return s.condominiumRepo.Create(condominium)
}

func (s *condominiumService) GetCondominiumByID(id uint, params map[string]interface{}) (*models.Condominium, error) {
	var condominium *models.Condominium
	var err error
	preload := params["preload"].(bool)
	log.Debug("preload: ", preload)

	if preload {
		condominium, err = s.condominiumRepo.FindByIDWithPreload(condominium, id, "Buildings")
		if err != nil {
			return nil, err
		}

	} else {
		condominium, err = s.condominiumRepo.FindByID(condominium, id)
		if err != nil {
			return nil, err
		}
	}
	return condominium, nil

}

func (s *condominiumService) GetAllCondominiums(page, pageSize int, params map[string]interface{}) ([]models.Condominium, int64, error) {
	log.Debug("Getting all condominiums")

	var condominiums []models.Condominium
	var totalRecords int64
	var err error
	preload := params["preload"].(bool)

	if preload {
		log.Debug("Preloading buildings")
		condominiums, totalRecords, err = s.condominiumRepo.FindAllWithPreload(page, pageSize)
		if err != nil {
			return nil, 0, err
		}
	} else {
		log.Debug("Not preloading buildings")
		condominiums, totalRecords, err = s.condominiumRepo.FindAll(page, pageSize)
		if err != nil {
			return nil, 0, err
		}
	}

	return condominiums, totalRecords, nil
}

func (s *condominiumService) UpdateCondominium(condominium *models.Condominium) error {
	var err error
	var condominiumAux *models.Condominium

	condominiumAux, err = s.condominiumRepo.FindByField(condominiumAux, "name", condominium.Name)
	if err != nil {
		return errors.New("erro to find condominium by name")
	}

	if condominiumAux != nil && condominium.ID != condominiumAux.ID {
		return errors.New("condominium already exists")
	}

	condominiumAux, err = s.condominiumRepo.FindByID(condominiumAux, condominium.ID)
	if err != nil {
		return errors.New("error to find condominium by id")
	}

	if condominiumAux == nil {
		return errors.New("condominium not found")
	}

	condominium.CreatedBy = condominiumAux.CreatedBy
	return s.condominiumRepo.Update(condominium)
}

func (s *condominiumService) DeleteCondominium(id uint) error {
	var condominium *models.Condominium
	condominium, err := s.condominiumRepo.FindByID(condominium, id)
	if err != nil {
		return err
	}
	return s.condominiumRepo.Delete(condominium)
}
