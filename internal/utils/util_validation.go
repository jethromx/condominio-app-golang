package utils

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct valida una estructura contra las etiquetas de validación definidas
func ValidateStruct(input interface{}) error {
	if err := validate.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationErrors(validationErrors)
		}
		log.Debug("Validation error: ", err)
		return errors.New("Invalid request " + err.Error())
	}
	return nil
}

// formatValidationErrors formatea los errores de validación en un mensaje detallado
func formatValidationErrors(validationErrors validator.ValidationErrors) error {
	var errorMessages []string
	for _, err := range validationErrors {
		errorMessage := fmt.Sprintf("Field '%s' failed validation with tag '%s'", err.Field(), err.Tag())
		errorMessages = append(errorMessages, errorMessage)
	}
	return errors.New(strings.Join(errorMessages, ", "))
}

// ValidEmail verifica si el formato del correo electrónico es válido
func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func HandlerValidation(c *fiber.Ctx, input interface{}) error {

	if err := c.BodyParser(input); err != nil {
		log.Debug("Error parsing input: ", err)
		return errors.New("invalid request, can't process paylod ")
	}

	if err := ValidateStruct(input); err != nil {
		log.Debug("Validation error: ", err)
		return errors.New("invalid request:  " + err.Error())
	}
	return nil
}
