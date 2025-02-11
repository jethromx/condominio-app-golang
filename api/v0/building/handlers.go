package building

import (
	"strconv"

	"com.mx/crud/api/v0/apartments"
	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// CreateBuildingHandler maneja la creación de nuevos edificios
func CreateBuildingHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(BuildingInput)

		err := utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		building := MapBuildingInputToModel(input)

		if err := buildingService.CreateBuilding(building); err != nil {
			log.Debug("Error creating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not create record, "+err.Error(), nil)
		}

		output := MapBuildingModelToOutput(building)
		return utils.HandleResponse(c, fiber.StatusOK, "Building created successfully", output)
	}
}

// GetBuildingByIDHandler maneja la obtención de un edificio por ID
func GetBuildingByIDHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			log.Debug("Error parsing ID:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid  ID ", nil)
		}

		building, err := buildingService.GetBuildingByID(uint(id))
		if err != nil {
			log.Debug("Error getting record:", err)
			return utils.HandleResponse(c, fiber.StatusNotFound, "Error: "+err.Error(), nil)
		}

		output := MapBuildingModelToOutput(building)

		return utils.HandleResponse(c, fiber.StatusOK, "Building retrieved successfully", output)
	}
}

// GetAllBuildingsHandler maneja la obtención de todos los edificios
func GetAllBuildingsHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", nil)
		}

		entities, totalRecords, err := buildingService.GetAllBuildings(page, pageSize)
		if err != nil {
			log.Debug("Error getting records:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error: "+err.Error(), nil)
		}

		output := MapBuildingsModelsToOutputs(entities)

		response := utils.DataPagination(page, pageSize, int(totalRecords), output)

		return utils.HandleResponse(c, fiber.StatusOK, "Users retrieved successfully", response)
	}
}

// UpdateBuildingHandler maneja la actualización de un edificio
func UpdateBuildingHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

		input := new(BuildingInput)
		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		building := MapBuildingInputToModel(input)
		building.ID = uint(id)

		if err := buildingService.UpdateBuilding(building); err != nil {
			log.Debug("Error updating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not update record, "+err.Error(), nil)
		}

		output := MapBuildingModelToOutput(building)
		return utils.HandleResponse(c, fiber.StatusOK, "Record updated successfully", output)
	}
}

// DeleteBuildingHandler maneja la eliminación de un edificio
func DeleteBuildingHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID "+err.Error(), nil)
		}

		if err := buildingService.DeleteBuilding(uint(id), 1); err != nil {
			log.Debug("Error deleting record:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not delete record "+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Building deleted successfully", nil)
	}
}

func GetAllApartmentsHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

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
			"id":      id,
		}
		log.Debug("Params: ", id)

		entities, totalRecords, err := apartmentService.GetAllApartments(page, pageSize, params)
		if err != nil {
			log.Debug("Error getting records:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error: "+err.Error(), nil)
		}

		output := apartments.MapApartmentsModelsToOutputs(entities)

		response := utils.DataPagination(page, pageSize, int(totalRecords), output)
		return utils.HandleResponse(c, fiber.StatusOK, "records retrieved successfully", response)
	}
}
func GetApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := c.ParamsInt("idApartment")
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

		output := apartments.MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, "record retrieved successfully", output)
	}
}
func CreateApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

		input := new(apartments.ApartmentInput)
		input.BuildingID = uint(id)

		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		apartment := apartments.MapAparmentInputToModel(input)

		if err := apartmentService.CreateApartment(apartment); err != nil {
			log.Debug("Error creating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not create record, "+err.Error(), nil)
		}

		output := apartments.MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, "Apartment created successfully", output)
	}
}
func UpdateApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("idApartment")
		input := new(apartments.ApartmentInput)

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

		idBuilding, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

		input.BuildingID = uint(idBuilding)
		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		apartment := apartments.MapAparmentInputToModel(input)
		apartment.ID = uint(id)

		if err := apartmentService.UpdateApartment(apartment); err != nil {
			log.Debug("Error updating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not update record, "+err.Error(), nil)
		}

		output := apartments.MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, "Record updated successfully", output)
	}

}
func DeleteApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("idApartment")
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
