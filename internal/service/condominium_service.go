package service

import (
	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
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

const (
	NAME      = "name"
	PRELOAD   = "preload"
	BUILDINGS = "Buildings"
	EMPTY     = ""
)

func NewCondominiumService(condominiumRepo repository.CondominiumRepository, buildingRepo repository.BuildingRepository) CondominiumService {
	return &condominiumService{condominiumRepo, buildingRepo}
}

func (s *condominiumService) CreateCondominium(condominium *models.Condominium) error {

	if condominium == nil {
		return models.ErrNilCondominium
	}

	var aux = &models.Condominium{}
	err := s.condominiumRepo.FindField(aux, NAME, condominium.Name)

	// Existe un condominio con el mismo nombre
	if err != nil {
		return models.ErrCondominiumNotFound
	}

	// Existe un condominio con el mismo nombre
	if aux.ID != 0 {
		return models.ErrCondominiumAlreadyExists
	}

	return s.condominiumRepo.Create(condominium)
}

func (s *condominiumService) GetCondominiumByID(id uint, params map[string]interface{}) (*models.Condominium, error) {
	var condominium = &models.Condominium{}
	var err error

	preload := params[PRELOAD].(bool)

	if preload {
		err = s.condominiumRepo.FindByIDPreload(condominium, id, BUILDINGS)
		if err != nil {
			return nil, err
		}

	} else {
		err = s.condominiumRepo.FindByIDPreload(condominium, id, EMPTY)
		if err != nil {
			return nil, err
		}
	}
	return condominium, nil

}

func (s *condominiumService) GetAllCondominiums(page, pageSize int, params map[string]interface{}) ([]models.Condominium, int64, error) {

	var condominiums []models.Condominium
	var totalRecords int64
	var err error
	preload := params[PRELOAD].(bool)

	if preload {
		condominiums, totalRecords, err = s.condominiumRepo.FindAllWithPreloadRel(page, pageSize, BUILDINGS)
		if err != nil {
			return nil, 0, err
		}
	} else {
		condominiums, totalRecords, err = s.condominiumRepo.FindAllWithPreloadRel(page, pageSize, EMPTY)
		if err != nil {
			return nil, 0, err
		}
	}

	return condominiums, totalRecords, nil
}

func (s *condominiumService) UpdateCondominium(condominium *models.Condominium) error {

	if err := s._ValidateCondominium(condominium); err != nil {
		return err
	}

	return s.condominiumRepo.Update(condominium)
}

func (s *condominiumService) DeleteCondominium(id uint) error {
	var condominium = &models.Condominium{}
	if err := s.condominiumRepo.FindID(condominium, id); err != nil {
		return err
	}

	return s.condominiumRepo.Delete(condominium)
}

func (s *condominiumService) _ValidateCondominium(condominium *models.Condominium) error {
	var condominiumAux = &models.Condominium{}
	var err error

	if err = s.condominiumRepo.FindField(condominiumAux, NAME, condominium.Name); err != nil {
		return models.ErrCondominiumNotFound
	}

	if condominiumAux.ID != 0 && condominium.ID != condominiumAux.ID {
		return models.ErrCondominiumAlreadyExists
	}

	if err = s.condominiumRepo.FindID(condominiumAux, condominium.ID); err != nil {
		return models.ErrCondominimByIdNotFound
	}

	if condominiumAux.ID == 0 {
		return models.ErrCondominiumNotFound
	}

	return nil
}
