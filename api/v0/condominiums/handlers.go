package condominiums

import (
	"strconv"

	"com.mx/crud/api/v0/apartments"
	"com.mx/crud/api/v0/building"
	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	AuditUserEMail = "email"
)

// @Router /condominiums [post]
func CreateCondominiumHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(CondominiumInput)

		if err := utils.HandlerValidation(c, input); err != nil {
			log.Debug(err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		auditUserID := c.Locals(AuditUserEMail).(string)

		input.CreatedBy = auditUserID
		input.UpdatedBy = auditUserID

		condominium := MapCondominiumInputToModel(input)

		if err := condominiumService.CreateCondominium(condominium); err != nil {
			log.Debug(utils.MsgErrorCreatingCondominium, err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorCreatingCondominium+err.Error(), nil)
		}

		output := MapCondominiumModelToOutput(condominium)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgCondominiumsRetrievedSuccessfully, output)
	}
}

func GetCondominiumByIDHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err error
		id, err := utils.GetParam(c, "id")
		if err != nil {
			log.Debug(utils.MsgErrorParsingID, err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		defaultParams := map[string]interface{}{
			"preload": false,
		}

		params, err := utils.GetQueryParams(c, defaultParams)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error getting parameters", nil)
		}

		condominium, err := condominiumService.GetCondominiumByID(uint(id), params)
		if err != nil {
			log.Debug(utils.MsgErrorGettingCondominiums, err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgErrorGettingCondominiums+err.Error(), nil)
		}

		output := MapCondominiumModelToOutput(condominium)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgCondominiumRetrievedSuccessfully, output)
	}
}

func GetAllCondominiumsHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		defaultParams := map[string]interface{}{
			"preload": false,
		}

		params, err := utils.GetQueryParams(c, defaultParams)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error getting parameters", nil)
		}

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		condominiums, totalRecords, err := condominiumService.GetAllCondominiums(page, pageSize, params)

		if err != nil {
			log.Debug(utils.MsgErrorGettingCondominiums, err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgErrorGettingCondominiums+err.Error(), nil)
		}

		output := MapCondominiumsModelsToOutputs(condominiums)
		response := utils.DataPagination(page, pageSize, int(totalRecords), output)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgCondominiumRetrievedSuccessfully, response)
	}
}

func UpdateCondominiumHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := utils.GetParam(c, "id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		input := new(CondominiumInput)

		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		auditUserID := c.Locals("email").(string)
		input.UpdatedBy = auditUserID

		condominium := MapCondominiumInputToModel(input)

		condominium.ID = uint(id) // Asegúrate de establecer el ID correcto

		if err := condominiumService.UpdateCondominium(condominium); err != nil {
			log.Debug(utils.MsgErrorUpdatingCondominium, err)
			log.Debug(err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorUpdatingCondominium+err.Error(), nil)
		}

		output := MapCondominiumModelToOutput(condominium)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgCondominiumUpdatedSuccessfully, output)
	}
}

func DeleteCondominiumHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := utils.GetParam(c, "id")

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		if err := condominiumService.DeleteCondominium(uint(id)); err != nil {
			log.Debug(utils.MsgErrorDeletingCondominium, err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorDeletingCondominium+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgCondominiumDeletedSuccessfully, nil)
	}
}

// ############################################################################################################

func GetAllBuildingsHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := utils.GetParam(c, "id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		defaultParams := map[string]interface{}{
			"preload": false,
		}

		params, err := utils.GetQueryParams(c, defaultParams)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error getting parameters", nil)
		}

		preload := params["preload"].(bool)

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", nil)
		}

		entities, totalRecords, err := buildingService.GetAllBuildingsByCondominium(page, pageSize, id, preload)
		if err != nil {
			log.Debug("Error getting records:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error: "+err.Error(), nil)
		}

		output := building.MapBuildingsModelsToOutputs(entities)

		response := utils.DataPagination(page, pageSize, int(totalRecords), output)

		return utils.HandleResponse(c, fiber.StatusOK, "Users retrieved successfully", response)
	}
}

func CreateBuildingsHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// GET ID FROM URL condominium
		id, err := utils.GetParam(c, "id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		input := new(building.BuildingInput)
		input.CondominiumID = uint(id)

		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		auditUserID := c.Locals("email").(string)
		input.UpdatedBy = auditUserID
		input.CreatedBy = auditUserID
		building := building.MapBuildingInputToModel(input)

		if err := buildingService.CreateBuilding(building); err != nil {
			log.Debug("Error creating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not create record, "+err.Error(), nil)
		}

		output := MapBuildingModelToOutput(building)
		return utils.HandleResponse(c, fiber.StatusOK, "Building created successfully", output)
	}
}

func GetBuildingsByIDHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("idBuilding")
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

func UpdateBuildingsHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idBuilding, err := c.ParamsInt("idBuilding")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

		// GET ID FROM URL condominium
		id, err := utils.GetParam(c, "id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		input := new(building.BuildingInput)
		input.CondominiumID = uint(id)
		auditUserID := c.Locals("email").(string)
		input.UpdatedBy = auditUserID

		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		building := building.MapBuildingInputToModel(input)
		building.ID = uint(idBuilding)

		if err := buildingService.UpdateBuilding(building); err != nil {
			log.Debug("Error updating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not update record, "+err.Error(), nil)
		}

		output := MapBuildingModelToOutput(building)
		return utils.HandleResponse(c, fiber.StatusOK, "Record updated successfully", output)
	}
}

func DeleteBuildingsHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idBuilding, err := c.ParamsInt("idBuilding")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

		// GET ID FROM URL condominium
		id, err := utils.GetParam(c, "id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		if err := buildingService.DeleteBuilding(uint(id), uint(idBuilding)); err != nil {
			log.Debug("Error deleting record:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not delete record "+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Building deleted successfully", nil)
	}
}

// ############################################################################################################

// CreateApartmentHandler handles the creation of new apartments
func CreateApartmentHandler(apartmentService service.ApartmentService, buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// GET ID FROM URL condominium
		id, err := utils.GetParam(c, "id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// GET ID FROM URL building
		idBuilding, err := utils.GetParam(c, "idBuilding")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		input := new(apartments.ApartmentInput)
		input.BuildingID = uint(idBuilding)

		// Validate input
		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		// Validate if building and condominium exist, return a condominium
		cond, err := apartmentService.ValidateApartment(idBuilding, id)

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Error validating condominium and building", nil)
		}

		if cond == nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "building or condominium does not exist", nil)
		}

		apartment := apartments.MapAparmentInputToModel(input)
		apartment.BuildingID = uint(idBuilding)

		if err := apartmentService.CreateApartment(apartment); err != nil {
			log.Debug("Error creating record: ", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not create record, "+err.Error(), nil)
		}

		output := apartments.MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, "Apartment created successfully", output)
	}
}

// GetApartmentHandler handles the retrieval of an apartment by ID

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

		id, err := c.ParamsInt("idBuilding")
		if err != nil {
			log.Debug("Error parsing ID:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID ", nil)
		}

		params := map[string]interface{}{
			"preload": preload,
			"id":      id,
		}

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

// UpdateApartmentHandler handles the update of an apartment
func UpdateApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("idApartment")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid ID", nil)
		}

		input := new(apartments.ApartmentInput)
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

// DeleteApartmentHandler handles the deletion of an apartment
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
