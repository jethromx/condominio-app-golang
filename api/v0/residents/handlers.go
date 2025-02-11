package residents

import (
	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// CreateResidentHandler maneja la creación de nuevos residentes
func CreateResidentHandler(residentService service.ResidentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(ResidentInput)

		err := utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		resident := MapResidentInputToModel(input)

		if err := residentService.CreateResident(resident); err != nil {
			log.Debug("Error creating resident: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not create record "+err.Error(), nil)
		}

		output := MapResidentToOutput(resident)

		return utils.HandleResponse(c, fiber.StatusOK, "Building created successfully", output)
	}
}

// GetResidentByIDHandler maneja la obtención de un residente por ID
func GetResidentByIDHandler(residentService service.ResidentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			log.Debug("Error parsing ID:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid condominium ID ", nil)
		}

		resident, err := residentService.GetResidentByID(uint(id))
		if err != nil {
			log.Debug("Error getting record:", err)
			return utils.HandleResponse(c, fiber.StatusNotFound, "Error: "+err.Error(), nil)
		}
		output := MapResidentToOutput(resident)

		return utils.HandleResponse(c, fiber.StatusOK, "Resident retrieved successfully", output)
	}
}

// GetAllResidentsHandler maneja la obtención de todos los residentes
func GetAllResidentsHandler(residentService service.ResidentService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", nil)
		}

		entities, totalRecords, err := residentService.GetAllResidents(page, pageSize)
		if err != nil {
			log.Debug("Error getting records:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error: "+err.Error(), nil)
		}

		output := MapResidentListToOutput(entities)

		response := utils.DataPagination(page, pageSize, int(totalRecords), output)

		return utils.HandleResponse(c, fiber.StatusOK, "Users retrieved successfully", response)
	}
}

// UpdateResidentHandler maneja la actualización de un residente
func UpdateResidentHandler(residentService service.ResidentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

		input := new(ResidentInput)
		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		resident := MapResidentInputToModel(input)
		resident.ID = uint(id)

		if err := residentService.UpdateResident(resident); err != nil {
			log.Debug("Error updating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not update record, "+err.Error(), nil)
		}

		output := MapResidentToOutput(resident)

		return utils.HandleResponse(c, fiber.StatusOK, "Record updated successfully", output)
	}
}

// DeleteResidentHandler maneja la eliminación de un residente
func DeleteResidentHandler(residentService service.ResidentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid resident ID", err.Error())
		}

		if err := residentService.DeleteResident(uint(id)); err != nil {
			log.Debug("Error deleting resident:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not delete record "+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Resident deleted successfully", nil)
	}
}
