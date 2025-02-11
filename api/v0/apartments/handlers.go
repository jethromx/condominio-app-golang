package apartments

import (
	"strconv"

	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// CreateApartmentHandler handles the creation of new apartments
func CreateApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(ApartmentInput)

		err := utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		apartment := MapAparmentInputToModel(input)

		if err := apartmentService.CreateApartment(apartment); err != nil {
			log.Debug("Error creating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not create record, "+err.Error(), nil)
		}

		output := MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, "Apartment created successfully", output)
	}
}

// GetApartmentHandler handles the retrieval of an apartment by ID

func GetApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := c.ParamsInt("id")
		if err != nil {
			log.Debug("Error parsing ID:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID ", nil)
		}

		preload, err := strconv.ParseBool(c.Query("preload", "false"))
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error parsing preload", nil)
		}

		params := map[string]interface{}{
			"preload": preload,
		}

		apartment, err := apartmentService.GetApartmentByID(uint(id), params)
		if err != nil {
			log.Debug("Error getting record:", err)
			return utils.HandleResponse(c, fiber.StatusNotFound, "Error: "+err.Error(), nil)
		}

		output := MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, "record retrieved successfully", output)
	}
}

func GetAllApartmentsHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", nil)
		}
		preload, err := strconv.ParseBool(c.Query("preload", "false"))
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error parsing preload", nil)
		}

		params := map[string]interface{}{
			"preload": preload,
		}

		entities, totalRecords, err := apartmentService.GetAllApartments(page, pageSize, params)
		if err != nil {
			log.Debug("Error getting records:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error: "+err.Error(), nil)
		}

		output := MapApartmentsModelsToOutputs(entities)

		response := utils.DataPagination(page, pageSize, int(totalRecords), output)
		return utils.HandleResponse(c, fiber.StatusOK, "records retrieved successfully", response)
	}
}

// UpdateApartmentHandler handles the update of an apartment
func UpdateApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

		input := new(ApartmentInput)
		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		apartment := MapAparmentInputToModel(input)
		apartment.ID = uint(id)

		if err := apartmentService.UpdateApartment(apartment); err != nil {
			log.Debug("Error updating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not update record, "+err.Error(), nil)
		}

		output := MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, "Record updated successfully", output)
	}
}

// DeleteApartmentHandler handles the deletion of an apartment
func DeleteApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID "+err.Error(), nil)
		}

		if err := apartmentService.DeleteApartment(uint(id)); err != nil {
			log.Debug("Error deleting record:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not delete record "+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, "record deleted successfully", nil)
	}
}
