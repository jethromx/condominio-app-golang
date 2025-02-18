package service

import (
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
	_ValidateCondominium(condominiumID uint, field string) error
	ValidateCondominiumAndBuilding(condominiumID uint, buildingID uint) (*models.Building, error)
}

const (
	APARTMENTS = "Apartments"
)

type buildingService struct {
	buildingRepo    repository.BuildingRepository
	condominiumRepo repository.CondominiumRepository
}

func NewBuildingService(buildingRepo repository.BuildingRepository, condominiumRepo repository.CondominiumRepository) BuildingService {
	return &buildingService{buildingRepo, condominiumRepo}
}

func (s *buildingService) CreateBuilding(building *models.Building) error {

	if building == nil {
		return models.ErrNilBuilding
	}

	if err := s._ValidateCondominium(building.CondominiumID, building.Name); err != nil {
		return err
	}

	return s.buildingRepo.Create(building)
}

func (s *buildingService) GetBuildingByID(id uint) (*models.Building, error) {
	var building = &models.Building{}
	if err := s.buildingRepo.FindByIDPreload(building, id, APARTMENTS); err != nil {
		return nil, err
	}
	return building, nil

}

func (s *buildingService) GetAllBuildings(page, pageSize int) ([]models.Building, int64, error) {
	return s.buildingRepo.FindAllWithPreloadRel(page, pageSize, APARTMENTS)
}

func (s *buildingService) UpdateBuilding(building *models.Building) error {
	var err error
	var entityAux = &models.Building{}

	err = s._ValidateCondominiumV2(building.CondominiumID, building.ID, building.Name)
	if err != nil {
		return err
	}

	if entityAux.ID != 0 && building.ID != entityAux.ID {
		return models.ErrBuildingAlreadyExists
	}

	//building.ID = entityAux.ID
	building.CreatedBy = entityAux.CreatedBy
	return s.buildingRepo.Update(building)
}

func (s *buildingService) DeleteBuilding(id uint, idBuilding uint) error {
	var building = &models.Building{}
	var err error

	var condominium = &models.Condominium{}
	if err := s.condominiumRepo.FindID(condominium, id); err != nil {
		return models.ErrCondominiumNotFound
	}

	if condominium.ID == 0 {
		return models.ErrBuildingNotFound
	}

	if err = s.buildingRepo.FindID(building, idBuilding); err != nil {
		return models.ErrBuildingByIdNotFound
	}

	if building.ID == 0 {
		return models.ErrBuildingNotFound
	}

	return s.buildingRepo.Delete(building)
}

func (s *buildingService) _ValidateCondominiumV2(condominiumID uint, buildingID uint, field string) error {
	var condominium = &models.Condominium{}
	var building = &models.Building{}
	var err error

	if err = s.condominiumRepo.FindID(condominium, condominiumID); err != nil {
		return models.ErrCondominiumNotFound
	}

	if err = s.buildingRepo.FindID(building, buildingID); err != nil {
		return models.ErrBuildingByIdNotFound
	}

	if condominium.ID == 0 {
		return models.ErrCondominiumNotFound
	} else {
		building, err := s.buildingRepo.ValidateBuilding(int(condominiumID), field)
		if err != nil {
			return models.ErrBuildingNameNotFound
		}
		if building != nil {
			if building.ID != buildingID {
				return models.ErrBuldingSameName
			} else {
				return nil
			}

		}
	}

	return nil

}

func (s *buildingService) _ValidateCondominium(condominiumID uint, field string) error {
	var condominium = &models.Condominium{}
	var err error

	if err = s.condominiumRepo.FindID(condominium, condominiumID); err != nil {
		return models.ErrCondominiumNotFound
	}

	if condominium.ID == 0 {
		return models.ErrCondominiumNotFound
	} else {
		building, err := s.buildingRepo.ValidateBuilding(int(condominiumID), field)
		if err != nil {
			log.Debug("no se encontro building con el mismo nombre")
			return models.ErrBuildingNameNotFound
		}
		if building != nil {
			log.Debug("ya existe un building con el mismo nombre")
			return models.ErrBuldingSameName
		}
	}

	return nil

}

func (s *buildingService) ValidateCondominiumAndBuilding(condominiumID uint, buildingID uint) (*models.Building, error) {

	building, err := s.buildingRepo.ValidateBuildingByIDs(int(condominiumID), int(buildingID))
	if err != nil {
		return nil, models.ErrBuildingByIdNotFound
	}

	return building, nil

}

func (s *buildingService) GetAllBuildingsByCondominium(page, pageSize int, id int, preload bool) ([]models.Building, int64, error) {
	if preload {
		return s.buildingRepo.FindBuildingsByCondominiumPreload(page, pageSize, id, APARTMENTS)
	} else {
		return s.buildingRepo.FindBuildingsByCondominiumPreload(page, pageSize, id, EMPTY)
	}

}
