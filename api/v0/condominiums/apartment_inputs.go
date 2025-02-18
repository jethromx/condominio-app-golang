package condominiums

import "com.mx/crud/internal/models"

type ApartmentInput struct {
	Number     string          `json:"number"`
	Name       string          `json:"name"`
	Floor      int             `json:"floor"`
	BuildingID uint            `json:"building_id"`
	CreatedBy  string          `json:"created_by"`
	UpdatedBy  string          `json:"updated_by"`
	Residents  []ResidentInput `json:"residents"`
}

func MapResidentInputsToModels(inputs []ResidentInput) []models.Resident {
	residents := make([]models.Resident, len(inputs))
	for i, input := range inputs {
		residents[i] = models.Resident{
			ID: input.ID,
		}
	}
	return residents
}

type ApartmentOutput struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	Number     string          `json:"number"`
	Floor      int             `json:"floor"`
	BuildingID uint            `json:"building_id"`
	Residents  []ResidentInput `json:"residents"`
	UpdatedAt  string          `json:"updated_at"`
	CreatedAt  string          `json:"created_at"`
}

func MapAparmentInputToModel(input *ApartmentInput) *models.Apartment {
	return &models.Apartment{
		Name:       input.Name,
		Number:     input.Number,
		Floor:      input.Floor,
		BuildingID: input.BuildingID,

		Residents: MapResidentInputsToModels(input.Residents),
	}
}

func MapResidentModelToInput(resident *models.Resident) *ResidentInput {
	return &ResidentInput{
		ID: resident.ID,
	}
}

func MapResidentsModelsToInputs(residents []models.Resident) []ResidentInput {
	if residents == nil {
		return []ResidentInput{}
	}
	inputs := make([]ResidentInput, 0)
	if len(residents) == 0 {
		inputs := make([]ResidentInput, len(residents))
		for i, resident := range residents {
			inputs[i] = *MapResidentModelToInput(&resident)
		}
	}

	return inputs
}

func MapApartmentModelToOutput(apartment *models.Apartment) *ApartmentOutput {
	return &ApartmentOutput{
		ID:         apartment.ID,
		Number:     apartment.Number,
		Name:       apartment.Name,
		Floor:      apartment.Floor,
		BuildingID: apartment.BuildingID,
		Residents:  MapResidentsModelsToInputs(apartment.Residents),
		CreatedAt:  apartment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  apartment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func MapApartmentsModelsToOutputs(apartments []models.Apartment) []ApartmentOutput {
	outputs := make([]ApartmentOutput, len(apartments))
	for i, apartment := range apartments {
		outputs[i] = *MapApartmentModelToOutput(&apartment)
	}
	return outputs
}

// CreateApartmentHandler handles the creation of new apartments
