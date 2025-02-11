package maintenances

import (
	"time"

	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type MaintenanceInput struct {
	Description   string    `json:"description"`
	Date          time.Time `json:"date"`
	Cost          float64   `json:"cost"`
	CondominiumID uint      `json:"condominium_id"`
}

// CreateMaintenanceHandler maneja la creación de nuevos mantenimientos
func CreateMaintenanceHandler(maintenanceService service.MaintenanceService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(MaintenanceInput)

		if err := c.BodyParser(input); err != nil {
			log.Debug("Error parsing input:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid request", err.Error())
		}

		maintenance, err := maintenanceService.CreateMaintenance(input.Description, input.Date, input.Cost, input.CondominiumID)
		if err != nil {
			log.Debug("Error creating maintenance:", err)
			return utils.HandleResponse(c, fiber.StatusUnprocessableEntity, "Could not create maintenance", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Maintenance created successfully", maintenance)
	}
}

// GetMaintenanceByIDHandler maneja la obtención de un mantenimiento por ID
func GetMaintenanceByIDHandler(maintenanceService service.MaintenanceService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid maintenance ID", err.Error())
		}

		maintenance, err := maintenanceService.GetMaintenanceByID(uint(id))
		if err != nil {
			log.Debug("Error getting maintenance:", err)
			return utils.HandleResponse(c, fiber.StatusNotFound, "Maintenance not found", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Maintenance retrieved successfully", maintenance)
	}
}

// GetAllMaintenancesHandler maneja la obtención de todos los mantenimientos
func GetAllMaintenancesHandler(maintenanceService service.MaintenanceService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		maintenances, err := maintenanceService.GetAllMaintenances()
		if err != nil {
			log.Debug("Error getting maintenances:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not get maintenances", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Maintenances retrieved successfully", maintenances)
	}
}

// UpdateMaintenanceHandler maneja la actualización de un mantenimiento
func UpdateMaintenanceHandler(maintenanceService service.MaintenanceService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid maintenance ID", err.Error())
		}

		input := new(MaintenanceInput)
		if err := c.BodyParser(input); err != nil {
			log.Debug("Error parsing input:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid request", err.Error())
		}

		maintenance, err := maintenanceService.UpdateMaintenance(uint(id), input.Description, input.Date, input.Cost, input.CondominiumID)
		if err != nil {
			log.Debug("Error updating maintenance:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not update maintenance", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Maintenance updated successfully", maintenance)
	}
}

// DeleteMaintenanceHandler maneja la eliminación de un mantenimiento
func DeleteMaintenanceHandler(maintenanceService service.MaintenanceService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid maintenance ID", err.Error())
		}

		if err := maintenanceService.DeleteMaintenance(uint(id)); err != nil {
			log.Debug("Error deleting maintenance:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not delete maintenance", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Maintenance deleted successfully", nil)
	}
}
