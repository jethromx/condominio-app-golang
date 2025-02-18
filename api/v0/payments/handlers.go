package payments

import (
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
	ID_PAYMENT     = "idPayment"
	ID_RESIDENT    = "idResident"
)

// CreatePaymentHandler maneja la creación de nuevos pagos
func CreatePaymentHandler(paymentService service.PaymentService, residentService service.ResidentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(PaymentInput)

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
		err = residentService.ValidateApartmentSU(idBuilding, idCondominium, idApartment)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Error "+err.Error(), nil)
		}
		payment := MapPaymentInputToModel(input)
		err = paymentService.CreatePayment(payment)
		if err != nil {
			log.Debug("Error creating payment:", err)
			return utils.HandleResponse(c, fiber.StatusUnprocessableEntity, "Could not create payment", err.Error())
		}

		output := MapPaymentToOutput(payment)

		return utils.HandleResponse(c, fiber.StatusOK, "Payment created successfully", output)
	}
}

// GetPaymentByIDHandler maneja la obtención de un pago por ID
func GetPaymentByIDHandler(paymentService service.PaymentService, residentService service.ResidentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		idPayment, err := utils.GetParam(c, ID_PAYMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// Validate apartment
		err = residentService.ValidateApartmentSU(idBuilding, idCondominium, idApartment)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Error "+err.Error(), nil)
		}

		defaultParams := map[string]interface{}{
			PRELOAD: false,
		}

		params, err := utils.GetQueryParams(c, defaultParams)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, utils.MsgErrorGettingsParams, nil)
		}

		payment, err := paymentService.GetPaymentByID(uint(idPayment), params)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusNotFound, "Payment not found", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Payment retrieved successfully", payment)
	}
}

// GetAllPaymentsHandler maneja la obtención de todos los pagos
func GetAllPaymentsHandler(paymentService service.PaymentService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		page, pageSize, err := utils.GetPaginationParams(c)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", nil)
		}

		payments, _, err := paymentService.GetAllPayments(page, pageSize, 1 /*id*/, false)
		if err != nil {
			log.Debug("Error getting payments:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not get payments", err.Error())
		}

		output := MapPaymentToOutputList(payments)

		return utils.HandleResponse(c, fiber.StatusOK, "Payments retrieved successfully", output)
	}
}

// UpdatePaymentHandler maneja la actualización de un pago
func UpdatePaymentHandler(paymentService service.PaymentService, residentService service.ResidentService) fiber.Handler {
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

		idPayment, err := utils.GetParam(c, ID_PAYMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// Validate apartment
		err = residentService.ValidateApartmentSU(idBuilding, idCondominium, idApartment)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Error "+err.Error(), nil)
		}

		input := new(PaymentInput)
		err = utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		input.ID = uint(idPayment)
		payment := MapPaymentInputToModel(input)

		//input.ResidentID
		err = paymentService.UpdatePayment(payment)
		if err != nil {
			log.Debug("Error updating payment:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not update payment", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Payment updated successfully", payment)
	}
}

// DeletePaymentHandler maneja la eliminación de un pago
func DeletePaymentHandler(paymentService service.PaymentService, residentService service.ResidentService) fiber.Handler {
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

		idPayment, err := utils.GetParam(c, ID_PAYMENT)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, utils.MsgErrorParsingID, err.Error())
		}

		// Validate apartment
		err = residentService.ValidateApartmentSU(idBuilding, idCondominium, idApartment)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Error "+err.Error(), nil)
		}
		if err := paymentService.DeletePayment(uint(idPayment)); err != nil {
			log.Debug("Error deleting payment:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not delete payment", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Payment deleted successfully", nil)
	}
}
