package condominiums

import (
	"com.mx/crud/api/v0/building"
	"com.mx/crud/internal/models"
)

type BuildingOutput struct {
	ID            uint   `json:"id"`
	IDCondominium uint   `json:"id_condominium"`
	Name          string `json:"name"`
	Floors        int    `json:"floors"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type CondominiumInput struct {
	Name          string                   `json:"name" validate:"required,min=3,max=100"`
	Address       string                   `json:"address" validate:"required,min=10,max=255"`
	Phone         string                   `json:"phone" validate:"omitempty,min=10,max=15"`
	Email         string                   `json:"email" validate:"omitempty,email"`
	ZIPCode       string                   `json:"zip_code" validate:"omitempty,min=5,max=10"`
	CreatedBy     string                   `json:"created_by" validate:"omitempty,min=3,max=64"`
	UpdatedBy     string                   `json:"updated_by" validate:"omitempty,min=3,max=64"`
	BuildingInput []building.BuildingInput `json:"buildings" validate:"omitempty"`
}

type CondominiumOutput struct {
	ID             uint             `json:"id"`
	Name           string           `json:"name"`
	Address        string           `json:"address"`
	Phone          string           `json:"phone"`
	Email          string           `json:"email"`
	ZIPCode        string           `json:"zip_code"`
	CreatedBy      string           `json:"created_by"`
	UpdatedBy      string           `json:"updated_by"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
	BuildingOutput []BuildingOutput `json:"buildings"`
}

// MapCondominiumInputToModel mapea el input de Condominium a un modelo de Condominium
func MapCondominiumInputToModel(input *CondominiumInput) *models.Condominium {
	return &models.Condominium{
		Name:      input.Name,
		Address:   input.Address,
		Phone:     input.Phone,
		Email:     input.Email,
		ZIPCode:   input.ZIPCode,
		CreatedBy: input.CreatedBy,
		UpdatedBy: input.UpdatedBy,
	}
}

func MapCondominiumsModelsToOutputs(condominiums []models.Condominium) []CondominiumOutput {
	outputs := make([]CondominiumOutput, len(condominiums))
	for i, condominium := range condominiums {
		outputs[i] = *MapCondominiumModelToOutput(&condominium)
	}
	return outputs
}

func MapCondominiumModelToOutput(condominium *models.Condominium) *CondominiumOutput {
	return &CondominiumOutput{
		ID: condominium.ID,

		Name:           condominium.Name,
		Address:        condominium.Address,
		Phone:          condominium.Phone,
		Email:          condominium.Email,
		ZIPCode:        condominium.ZIPCode,
		CreatedBy:      condominium.CreatedBy,
		UpdatedBy:      condominium.UpdatedBy,
		CreatedAt:      condominium.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      condominium.UpdatedAt.Format("2006-01-02 15:04:05"),
		BuildingOutput: MapBuildingModelToOutputs(condominium.Buildings),
	}
}

func MapBuildingModelToOutputs(buildings []models.Building) []BuildingOutput {
	outputs := make([]BuildingOutput, len(buildings))
	for i, building := range buildings {
		outputs[i] = *MapBuildingModelToOutput(&building)
	}
	return outputs
}

func MapBuildingModelToOutput(building *models.Building) *BuildingOutput {

	return &BuildingOutput{
		ID:            building.ID,
		IDCondominium: building.CondominiumID,
		Name:          building.Name,
		Floors:        building.Floors,
		CreatedAt:     building.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     building.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
