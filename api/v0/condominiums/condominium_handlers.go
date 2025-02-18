package condominiums

import (
	"strconv"

	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	AuditUserEMail = "email"
	ID_CONDOMINIUM = "id"
	CREATE         = "create"
	PRELOAD        = "preload"
	ID_BUILDING    = "idBuilding"
	ID_APARTMENT   = "idApartment"
)

// @Router /condominiums [post]
func CreateCondominiumHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(CondominiumInput)

		if err := utils.HandlerValidation(c, input); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		auditUserID, err := utils.GetAuditUserID(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		setAuditUserID(auditUserID, input, CREATE)

		condominium := MapCondominiumInputToModel(input)

		if err := condominiumService.CreateCondominium(condominium); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorCreatingCondominium+err.Error(), nil)
		}

		output := MapCondominiumModelToOutput(condominium)

		return utils.HandleResponse(c, fiber.StatusCreated, utils.MsgCondominiumsRetrievedSuccessfully, output)
	}
}

func GetCondominiumByIDHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err error
		id, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			log.Debug(utils.MsgErrorParsingID, err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		defaultParams := map[string]interface{}{
			PRELOAD: false,
		}

		params, err := utils.GetQueryParams(c, defaultParams)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgErrorGettingsParams, nil)
		}

		condominium, err := condominiumService.GetCondominiumByID(uint(id), params)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgErrorGettingCondominiums+err.Error(), nil)
		}

		output := MapCondominiumModelToOutput(condominium)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgCondominiumRetrievedSuccessfully, output)
	}
}

func GetAllCondominiumsHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		defaultParams := map[string]interface{}{
			PRELOAD: false,
		}

		params, err := utils.GetQueryParams(c, defaultParams)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgErrorGettingsParams, nil)
		}

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		condominiums, totalRecords, err := condominiumService.GetAllCondominiums(page, pageSize, params)

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgErrorGettingCondominiums+err.Error(), nil)
		}

		output := MapCondominiumsModelsToOutputs(condominiums)
		response := utils.DataPagination(page, pageSize, int(totalRecords), output)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgCondominiumRetrievedSuccessfully, response)
	}
}

func UpdateCondominiumHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		input := new(CondominiumInput)

		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		auditUserID, err := utils.GetAuditUserID(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		setAuditUserID(auditUserID, input, CREATE)

		condominium := MapCondominiumInputToModel(input)
		condominium.ID = uint(id) // Asegúrate de establecer el ID correcto

		if err := condominiumService.UpdateCondominium(condominium); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorUpdatingCondominium+err.Error(), nil)
		}

		output := MapCondominiumModelToOutput(condominium)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgCondominiumUpdatedSuccessfully, output)
	}
}

func DeleteCondominiumHandler(condominiumService service.CondominiumService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := utils.GetParam(c, ID_CONDOMINIUM)

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		if err := condominiumService.DeleteCondominium(uint(id)); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorDeletingCondominium+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgCondominiumDeletedSuccessfully, nil)
	}
}

// ############################################################################################################

func GetAllBuildingsHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		defaultParams := map[string]interface{}{
			PRELOAD: false,
		}

		params, err := utils.GetQueryParams(c, defaultParams)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgErrorGettingsParams, nil)
		}

		preload := params[PRELOAD].(bool)

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgInvalidaParams, nil)
		}

		entities, totalRecords, err := buildingService.GetAllBuildingsByCondominium(page, pageSize, id, preload)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error: "+err.Error(), nil)
		}

		output := MapBuildingsModelsToOutputs(entities)

		response := utils.DataPagination(page, pageSize, int(totalRecords), output)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgBuildingGetSuccessfully, response)
	}
}

func CreateBuildingsHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// GET ID FROM URL condominium
		id, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		input := new(BuildingInput)
		input.CondominiumID = uint(id)

		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		auditUserID, err := utils.GetAuditUserID(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		setAuditUserBuildingID(auditUserID, input, CREATE)

		building := MapBuildingInputToModel(input)

		if err := buildingService.CreateBuilding(building); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgCouldGetRecord+err.Error(), nil)
		}
		log.Debug(building.ID)

		output := MapBuildingModelToOutput(building)
		return utils.HandleResponse(c, fiber.StatusCreated, utils.MsgBuildingCreatedSuccessfully, output)
	}
}

func GetBuildingsByIDHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt(ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		building, err := buildingService.GetBuildingByID(uint(id))
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusNotFound, "Error: "+err.Error(), nil)
		}

		output := MapBuildingModelToOutput(building)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgBuildingGetSuccessfully, output)
	}
}

func UpdateBuildingsHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idBuilding, err := c.ParamsInt(ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgInvalidaParams, nil)
		}

		// GET ID FROM URL condominium
		id, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		input := new(BuildingInput)
		input.CondominiumID = uint(id)
		auditUserID, err := utils.GetAuditUserID(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
		}
		input.UpdatedBy = auditUserID

		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		building := MapBuildingInputToModel(input)
		building.ID = uint(idBuilding)

		if err := buildingService.UpdateBuilding(building); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgCouldGetRecord+err.Error(), nil)
		}

		output := MapBuildingModelToOutput(building)
		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgBuildingUpdatedSuccessfully, output)
	}
}

func DeleteBuildingsHandler(buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idBuilding, err := c.ParamsInt(ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, nil)
		}

		// GET ID FROM URL condominium
		id, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		if err := buildingService.DeleteBuilding(uint(id), uint(idBuilding)); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgCouldGetRecord+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgBuildingDeletedSuccessfully, nil)
	}
}

// ############################################################################################################

// CreateApartmentHandler handles the creation of new apartments
func CreateApartmentHandler(apartmentService service.ApartmentService, buildingService service.BuildingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(ApartmentInput)
		if err := utils.HandlerValidation(c, input); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		// GET ID FROM URL condominium
		id, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// GET ID FROM URL building
		idBuilding, err := utils.GetParam(c, ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		auditUserID, err := utils.GetAuditUserID(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		setAuditUserApartmentID(auditUserID, input, CREATE)

		input.BuildingID = uint(idBuilding)

		// Validate if building and condominium exist, return a condominium
		cond, err := apartmentService.ValidateApartment(idBuilding, id)

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorValidationApartment, nil)
		}

		if cond == nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgCondominiumOrBuildingDoesNotExist, nil)
		}

		apartment := MapAparmentInputToModel(input)

		if len(apartment.Residents) > 0 {
			ids := make([]int, len(apartment.Residents))
			for i, resident := range apartment.Residents {
				ids[i] = int(resident.ID)
			}
			if err = apartmentService.ValidateResidents(ids); err != nil {
				return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgResidentsNotFound, nil)
			}
		}

		if err := apartmentService.CreateApartment(apartment); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgCouldNotCreateApartment+err.Error(), nil)
		}

		output := MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgApartmentCreatedSuccessfully, output)
	}
}

// GetApartmentHandler handles the retrieval of an apartment by ID

func GetApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := c.ParamsInt(ID_APARTMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		preload, err := strconv.ParseBool(c.Query(PRELOAD, "false"))

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgInvalidaParams, nil)
		}

		params := map[string]interface{}{
			PRELOAD: preload,
		}

		apartment, err := apartmentService.GetApartmentByID(uint(id), params)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusNotFound, "Error: "+err.Error(), nil)
		}

		if apartment == nil {
			return utils.HandleResponse(c, fiber.StatusNotFound, "Record not found", nil)
		}
		output := MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgApartmentsRetrievedSuccessfully, output)
	}
}

func GetAllApartmentsHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgInvalidaParams, nil)
		}
		preload, err := strconv.ParseBool(c.Query(PRELOAD, "false"))
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgInvalidaParams, nil)
		}

		id, err := c.ParamsInt(ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		params := map[string]interface{}{
			PRELOAD: preload,
			"id":    id,
		}

		entities, totalRecords, err := apartmentService.GetAllApartments(page, pageSize, params)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error: "+err.Error(), nil)
		}

		output := MapApartmentsModelsToOutputs(entities)

		response := utils.DataPagination(page, pageSize, int(totalRecords), output)
		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgApartmentsRetrievedSuccessfully, response)
	}
}

func GetAllApartmentsByCondominiumHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		condominiumID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid condominium ID",
			})
		}

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgInvalidaParams, nil)
		}
		/*
			preload, err := strconv.ParseBool(c.Query(PRELOAD, "false"))
			if err != nil {
				return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgInvalidaParams, nil)
			}*/

		apartments, totalRecords, err := apartmentService.GetAllApartmentsByCondominium(page, pageSize, uint(condominiumID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get apartments",
			})
		}

		response := utils.DataPagination(page, pageSize, int(totalRecords), apartments)

		return c.JSON(response)
	}
}

// UpdateApartmentHandler handles the update of an apartment
func UpdateApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt(ID_APARTMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		idCon, err := c.ParamsInt(ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// GET ID FROM URL building
		idBuilding, err := utils.GetParam(c, ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		input := new(ApartmentInput)
		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		// Validate if building and condominium exist, return a condominium
		cond, err := apartmentService.ValidateApartment(idBuilding, idCon)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorValidationApartment, nil)
		}

		if cond == nil || cond.ID != uint(idCon) {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgCondominiumOrBuildingDoesNotExist, nil)
		}

		apartment := MapAparmentInputToModel(input)
		apartment.ID = uint(id)
		apartment.BuildingID = uint(idBuilding)

		if len(apartment.Residents) > 0 {
			ids := make([]int, len(apartment.Residents))
			for i, resident := range apartment.Residents {
				ids[i] = int(resident.ID)
			}
			if err = apartmentService.ValidateResidents(ids); err != nil {
				return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgResidentsNotFound, nil)
			}
		}

		if err := apartmentService.UpdateApartment(apartment); err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Could not update record, "+err.Error(), nil)
		}

		params := map[string]interface{}{
			PRELOAD: true,
		}

		apartment, err = apartmentService.GetApartmentByID(apartment.ID, params)

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error: "+err.Error(), nil)
		}

		output := MapApartmentModelToOutput(apartment)

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgApartmentUpdatedSuccessfully, output)
	}
}

// DeleteApartmentHandler handles the deletion of an apartment
func DeleteApartmentHandler(apartmentService service.ApartmentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt(ID_APARTMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		idCon, err := c.ParamsInt(ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// GET ID FROM URL building
		idBuilding, err := utils.GetParam(c, ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// Validate if building and condominium exist, return a condominium
		cond, err := apartmentService.ValidateApartment(idBuilding, idCon)

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorValidationApartment, nil)
		}

		if cond == nil || cond.ID != uint(idCon) {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgCondominiumOrBuildingDoesNotExist, nil)
		}

		if err := apartmentService.DeleteApartment(uint(id)); err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Error "+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, utils.MsgApartmentDeletedSuccessfully, nil)
	}
}
