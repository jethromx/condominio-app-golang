package utils

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
)

// GetPaginationParams extrae y valida los parámetros de paginación de la solicitud
func GetPaginationParams(c *fiber.Ctx) (int, int, error) {
	page, err := strconv.Atoi(c.Query("page", strconv.Itoa(defaultPage)))
	if err != nil || page < 1 {
		page = defaultPage
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", strconv.Itoa(defaultPageSize)))
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}

	return page, pageSize, nil
}

// GetParams obtiene los parámetros de la consulta y los combina con los parámetros predeterminados
func GetQueryParams(c *fiber.Ctx, defaultParams map[string]interface{}) (map[string]interface{}, error) {
	params := make(map[string]interface{})

	// Copiar los parámetros predeterminados al mapa de parámetros
	for key, value := range defaultParams {
		params[key] = value
	}

	// Obtener los parámetros de la consulta
	for key := range c.Queries() {
		value := c.Query(key)
		// Intentar convertir el valor a booleano
		if boolValue, err := strconv.ParseBool(value); err == nil {
			params[key] = boolValue
		} else {
			// Si no es un booleano, mantenerlo como cadena
			params[key] = value
		}
	}

	return params, nil
}

func GetParam(c *fiber.Ctx, paramId string) (int, error) {
	param, err := c.ParamsInt(paramId)
	if err != nil {
		log.Debug("Error parsing param: ", err)
		return 0, errors.New("Invalid request " + err.Error())
	}
	return param, nil

}
