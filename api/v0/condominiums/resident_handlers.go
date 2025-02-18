package condominiums

import (
	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var (
	ID_RESIDENT = "idResident"
)

// CreateResidentHandler maneja la creación de nuevos residentes
func CreateResidentHandler(residentService service.ResidentService, userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(ResidentInput)

		// GET ID FROM URL condominium
		idCondominium, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// GET ID FROM URL building
		idBuilding, err := utils.GetParam(c, ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		idApartment, err := utils.GetParam(c, ID_APARTMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		// Validate apartment
		err = residentService.ValidateApartment(idBuilding, idCondominium, idApartment, int(input.UserID))
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Error "+err.Error(), nil)
		}

		input.ApartmentID = uint(idApartment)
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
		// GET ID FROM URL condominium
		idCondominium, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// GET ID FROM URL building
		idBuilding, err := utils.GetParam(c, ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		idApartment, err := utils.GetParam(c, ID_APARTMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		idResident, err := utils.GetParam(c, ID_RESIDENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// Validate apartment
		err = residentService.ValidateApartmentSU(idBuilding, idCondominium, idApartment)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Error "+err.Error(), nil)
		}

		resident, err := residentService.GetResidentByID(uint(idResident))
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
		// GET ID FROM URL condominium
		idCondominium, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// GET ID FROM URL building
		idBuilding, err := utils.GetParam(c, ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		idApartment, err := utils.GetParam(c, ID_APARTMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		idResident, err := utils.GetParam(c, ID_RESIDENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// Validate apartment
		err = residentService.ValidateApartmentSU(idBuilding, idCondominium, idApartment)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Error "+err.Error(), nil)
		}

		input := new(ResidentInput)
		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		resident := MapResidentInputToModel(input)
		resident.ID = uint(idResident)
		resident.ApartmentID = uint(idApartment)

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
		// GET ID FROM URL condominium
		idCondominium, err := utils.GetParam(c, ID_CONDOMINIUM)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// GET ID FROM URL building
		idBuilding, err := utils.GetParam(c, ID_BUILDING)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		idApartment, err := utils.GetParam(c, ID_APARTMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		idResident, err := utils.GetParam(c, ID_RESIDENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// Validate apartment
		err = residentService.ValidateApartmentSU(idBuilding, idCondominium, idApartment)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Error "+err.Error(), nil)
		}

		if err := residentService.DeleteResident(uint(idResident)); err != nil {
			log.Debug("Error deleting resident:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not delete record "+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Resident deleted successfully", nil)
	}
}
