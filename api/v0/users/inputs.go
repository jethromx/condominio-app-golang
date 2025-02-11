package users

import "com.mx/crud/internal/models"

type UserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"rol" validate:"omitempty,oneof=admin manager resident"`
}

type UserOutput struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Role      string `json:"rol" validate:"omitempty,oneof=admin manager resident"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserOutputs struct {
	Users []*UserOutput `json:"users"`
}

// MapUserInputToModel mapea el input de User a un modelo de User
func MapUserInputToModel(input *UserInput) *models.User {
	return &models.User{
		Email:    input.Email,
		Username: input.Username,
		Password: input.Password,
		Role:     input.Role,
	}
}

// MapUserModelToOutput mapea el modelo de User a un output de User
func MapUserModelToOutput(user *models.User) *UserOutput {
	return &UserOutput{
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      user.Role,
	}
}

// MapUserModelsToOutputs mapea un slice de modelos de User a un slice de outputs de User
func MapUserModelsToOutputs(users []models.User) []UserOutput {
	outputs := make([]UserOutput, len(users))
	for i, user := range users {
		outputs[i] = *MapUserModelToOutput(&user)
	}
	return outputs
}
