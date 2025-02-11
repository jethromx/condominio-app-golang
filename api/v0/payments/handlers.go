package payments

import (
	"time"

	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type PaymentInput struct {
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
	ResidentID uint      `json:"resident_id"`
}

// CreatePaymentHandler maneja la creación de nuevos pagos
func CreatePaymentHandler(paymentService service.PaymentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(PaymentInput)

		if err := c.BodyParser(input); err != nil {
			log.Debug("Error parsing input:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid request", err.Error())
		}

		payment, err := paymentService.CreatePayment(input.Amount, input.Date, input.ResidentID)
		if err != nil {
			log.Debug("Error creating payment:", err)
			return utils.HandleResponse(c, fiber.StatusUnprocessableEntity, "Could not create payment", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Payment created successfully", payment)
	}
}

// GetPaymentByIDHandler maneja la obtención de un pago por ID
func GetPaymentByIDHandler(paymentService service.PaymentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid payment ID", err.Error())
		}

		payment, err := paymentService.GetPaymentByID(uint(id))
		if err != nil {
			log.Debug("Error getting payment:", err)
			return utils.HandleResponse(c, fiber.StatusNotFound, "Payment not found", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Payment retrieved successfully", payment)
	}
}

// GetAllPaymentsHandler maneja la obtención de todos los pagos
func GetAllPaymentsHandler(paymentService service.PaymentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payments, err := paymentService.GetAllPayments()
		if err != nil {
			log.Debug("Error getting payments:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not get payments", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Payments retrieved successfully", payments)
	}
}

// UpdatePaymentHandler maneja la actualización de un pago
func UpdatePaymentHandler(paymentService service.PaymentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid payment ID", err.Error())
		}

		input := new(PaymentInput)
		if err := c.BodyParser(input); err != nil {
			log.Debug("Error parsing input:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid request", err.Error())
		}

		payment, err := paymentService.UpdatePayment(uint(id), input.Amount, input.Date, input.ResidentID)
		if err != nil {
			log.Debug("Error updating payment:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not update payment", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Payment updated successfully", payment)
	}
}

// DeletePaymentHandler maneja la eliminación de un pago
func DeletePaymentHandler(paymentService service.PaymentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid payment ID", err.Error())
		}

		if err := paymentService.DeletePayment(uint(id)); err != nil {
			log.Debug("Error deleting payment:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not delete payment", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Payment deleted successfully", nil)
	}
}
