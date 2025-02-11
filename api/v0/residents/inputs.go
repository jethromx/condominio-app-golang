package residents

import "com.mx/crud/internal/models"

type ResidentInput struct {
	ApartmentID uint   `json:"apartment_id" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Phone       string `json:"phone" `
	Email       string `json:"email" validate:"required,email"`
	UserID      uint   `json:"user_id" validate:"required"`
}

type ResidentOutput struct {
	ID          uint   `json:"id"`
	ApartmentID uint   `json:"apartment_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	UserID      uint   `json:"user_id"`
}

func MapResidentInputToModel(input *ResidentInput) *models.Resident {

	return &models.Resident{
		ApartmentID: input.ApartmentID,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Phone:       input.Phone,
		Email:       input.Email,
		UserID:      input.UserID,
	}
}

func MapResidentToOutput(resident *models.Resident) *ResidentOutput {

	return &ResidentOutput{
		ID:          resident.ID,
		ApartmentID: resident.ApartmentID,
		FirstName:   resident.FirstName,
		LastName:    resident.LastName,
		Phone:       resident.Phone,
		Email:       resident.Email,
		UserID:      resident.UserID,
	}
}

func MapResidentListToOutput(residents []models.Resident) []ResidentOutput {
	var output []ResidentOutput

	for _, resident := range residents {
		output = append(output, *MapResidentToOutput(&resident))
	}

	return output
}
