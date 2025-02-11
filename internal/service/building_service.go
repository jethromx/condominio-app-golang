package service

import (
	"errors"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
	"github.com/gofiber/fiber/v2/log"
)

type BuildingService interface {
	CreateBuilding(building *models.Building) error
	GetBuildingByID(id uint) (*models.Building, error)
	GetAllBuildings(page, pageSize int) ([]models.Building, int64, error)
	GetAllBuildingsByCondominium(page, pageSize int, id int, preload bool) ([]models.Building, int64, error)
	UpdateBuilding(building *models.Building) error
	DeleteBuilding(id uint, idBuilding uint) error
	ValidateCondominium(condominiumID uint, field string) error
	ValidateCondominiumAndBuilding(condominiumID uint, buildingID uint) (*models.Building, error)
}

type buildingService struct {
	buildingRepo    repository.BuildingRepository
	condominiumRepo repository.CondominiumRepository
}

func NewBuildingService(buildingRepo repository.BuildingRepository, condominiumRepo repository.CondominiumRepository) BuildingService {
	return &buildingService{buildingRepo, condominiumRepo}
}

func (s *buildingService) CreateBuilding(building *models.Building) error {

	if building == nil {
		log.Debug("Building entity is nil")
		return errors.New("error input data")

	}

	if err := s.ValidateCondominium(building.CondominiumID, building.Name); err != nil {
		return err
	}

	return s.buildingRepo.Create(building)
}

func (s *buildingService) GetBuildingByID(id uint) (*models.Building, error) {
	var building *models.Building
	return s.buildingRepo.FindByIDWithPreload(building, id, "Apartments")
}

func (s *buildingService) GetAllBuildings(page, pageSize int) ([]models.Building, int64, error) {
	return s.buildingRepo.FindAllWithPreloadRel(page, pageSize, "Apartments")
}

func (s *buildingService) UpdateBuilding(building *models.Building) error {
	var err error
	var entityAux *models.Building

	err = s.ValidateCondominium(building.CondominiumID, building.Name)
	if err != nil {
		return err
	}

	entityAux, err = s.buildingRepo.FindByID(entityAux, building.ID)
	if err != nil {
		return err
	}

	_, err = s.buildingRepo.FindByField(entityAux, "name", building.Name)

	if err != nil {
		return err
	}

	if entityAux != nil && building.ID != entityAux.ID {
		return errors.New("record already exists")
	}

	building.ID = entityAux.ID
	building.CreatedBy = entityAux.CreatedBy
	return s.buildingRepo.Update(building)
}

func (s *buildingService) DeleteBuilding(id uint, idBuilding uint) error {
	var building *models.Building

	var condominium *models.Condominium
	condominium, err := s.condominiumRepo.FindByID(condominium, id)

	if err != nil {
		log.Debug("Error finding condominium by ID ", err)
		return errors.New("error finding condominium by ID")
	}
	if condominium == nil {
		return errors.New("condominium does not exist")
	}

	building, err = s.buildingRepo.FindByID(building, idBuilding)

	if err != nil {
		log.Debug("Error finding building by ID ", err)
		return errors.New("error finding building by ID")
	}

	if building == nil {
		log.Debug("Building does not exist")
		return errors.New("building does not exist")
	}

	return s.buildingRepo.Delete(building)
}

func (s *buildingService) ValidateCondominium(condominiumID uint, field string) error {
	var condominium *models.Condominium
	condominium, err := s.condominiumRepo.FindByID(condominium, condominiumID)

	if err != nil {
		log.Debug("Error finding condominium by ID ", err)
		return errors.New("error finding condominium by ID")
	}

	if condominium == nil {
		return errors.New("condominium does not exist")
	} else {
		building, err := s.buildingRepo.ValidateBuilding(int(condominiumID), field)
		if err != nil {
			log.Debug("Error finding building by name", err)
			return errors.New("error finding building by name")
		}
		if building != nil {
			log.Debug("Building already exists")
			return errors.New("in this condominium a building already exists with the same name")
		}

		// No existe un condominio con el mismo nombre

	}

	return nil

}

func (s *buildingService) ValidateCondominiumAndBuilding(condominiumID uint, buildingID uint) (*models.Building, error) {

	building, err := s.buildingRepo.ValidateBuildingByIDs(int(condominiumID), int(buildingID))
	if err != nil {
		log.Debug("Error finding building by name", err)
		return nil, errors.New("error finding building by name")
	}
	/*if building != nil {
		log.Debug("Building already exists")
		return building, errors.New("in this condominium a building already exists with the same name")
	}*/
	return building, nil

}

func (s *buildingService) GetAllBuildingsByCondominium(page, pageSize int, id int, preload bool) ([]models.Building, int64, error) {
	if preload {
		return s.buildingRepo.FindBuildingsByCondominiumPreload(page, pageSize, id)
	} else {
		return s.buildingRepo.FindBuildingsByCondominium(page, pageSize, id)
	}

}
