package building

import "com.mx/crud/internal/models"

// BuildingInput is the input for creating a building
type BuildingInput struct {
	CondominiumID uint             `json:"id_condominium" validate:"omitempty"`
	Name          string           `json:"name" validate:"required,min=3,max=100"`
	Floors        int              `json:"floors" validate:"required,min=1"`
	CreatedBy     string           `json:"created_by" validate:"omitempty,min=3,max=100"`
	UpdatedBy     string           `json:"update_by" validate:"omitempty,min=3,max=100"`
	Apartments    []ApartmentInput `json:"apartments"`
}

type BuildingOutput struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name"`
	Floors       int            `json:"floors"`
	ApartmentOut []ApartmentOut `json:"apartments"`
	CreatedBy    string         `json:"created_by"`
	UpdateBy     string         `json:"update_by"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
}

type ApartmentOut struct {
	IDBuilding uint   `json:"id_building"`
	ID         uint   `json:"id"`
	Number     string `json:"number"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Floor      int    `json:"floor"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ApartmentInput is the input for creating an apartment
type ApartmentInput struct {
	Number int `json:"id_apartment" validate:"required"`
}

// MapBuildingInputToModel maps the BuildingInput to a Building model

func MapBuildingInputToModel(input *BuildingInput) *models.Building {

	return &models.Building{
		CondominiumID: input.CondominiumID,
		Name:          input.Name,
		Floors:        input.Floors,
		CreatedBy:     input.CreatedBy,
		UpdatedBy:     input.UpdatedBy,
	}
}

func MapBuildingModelToOutput(building *models.Building) *BuildingOutput {
	return &BuildingOutput{
		ID:           building.ID,
		Name:         building.Name,
		Floors:       building.Floors,
		ApartmentOut: MapApartmentOutputs(building.Apartments),
		CreatedAt:    building.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    building.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:    building.CreatedBy,
		UpdateBy:     building.UpdatedBy,
	}
}

func MapBuildingsModelsToOutputs(buildings []models.Building) []BuildingOutput {
	outputs := make([]BuildingOutput, len(buildings))
	for i, building := range buildings {
		outputs[i] = *MapBuildingModelToOutput(&building)
	}
	return outputs
}

func MapAparmentToOutput(apartment *models.Apartment) *ApartmentOut {
	return &ApartmentOut{
		ID:         apartment.ID,
		IDBuilding: apartment.BuildingID,
		Number:     apartment.Number,
		Name:       apartment.Name,
		Status:     apartment.Status,
		Floor:      apartment.Floor,
		CreatedAt:  apartment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  apartment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func MapApartmentOutputs(apartments []models.Apartment) []ApartmentOut {
	outputs := make([]ApartmentOut, len(apartments))
	for i, apartment := range apartments {
		outputs[i] = *MapAparmentToOutput(&apartment)
	}
	return outputs
}
