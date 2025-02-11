package auth

import "com.mx/crud/internal/models"

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginInput struct {
	Identity string `json:"identity" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func MapUserInputToModel(input *RegisterInput) *models.User {
	return &models.User{
		Email:    input.Email,
		Username: input.Username,
		Password: input.Password, // Asegúrate de manejar el hashing de la contraseña en otro lugar
	}
}
