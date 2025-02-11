package reservations

import (
	"time"

	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ReservationInput struct {
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	ResidentID uint      `json:"resident_id"`
}

// CreateReservationHandler maneja la creación de nuevas reservas
func CreateReservationHandler(reservationService service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(ReservationInput)

		if err := c.BodyParser(input); err != nil {
			log.Debug("Error parsing input:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid request", err.Error())
		}

		reservation, err := reservationService.CreateReservation(input.StartDate, input.EndDate, input.ResidentID)
		if err != nil {
			log.Debug("Error creating reservation:", err)
			return utils.HandleResponse(c, fiber.StatusUnprocessableEntity, "Could not create reservation", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Reservation created successfully", reservation)
	}
}

// GetReservationByIDHandler maneja la obtención de una reserva por ID
func GetReservationByIDHandler(reservationService service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid reservation ID", err.Error())
		}

		reservation, err := reservationService.GetReservationByID(uint(id))
		if err != nil {
			log.Debug("Error getting reservation:", err)
			return utils.HandleResponse(c, fiber.StatusNotFound, "Reservation not found", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Reservation retrieved successfully", reservation)
	}
}

// GetAllReservationsHandler maneja la obtención de todas las reservas
func GetAllReservationsHandler(reservationService service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reservations, err := reservationService.GetAllReservations()
		if err != nil {
			log.Debug("Error getting reservations:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not get reservations", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Reservations retrieved successfully", reservations)
	}
}

// UpdateReservationHandler maneja la actualización de una reserva
func UpdateReservationHandler(reservationService service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid reservation ID", err.Error())
		}

		input := new(ReservationInput)
		if err := c.BodyParser(input); err != nil {
			log.Debug("Error parsing input:", err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid request", err.Error())
		}

		reservation, err := reservationService.UpdateReservation(uint(id), input.StartDate, input.EndDate, input.ResidentID)
		if err != nil {
			log.Debug("Error updating reservation:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not update reservation", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Reservation updated successfully", reservation)
	}
}

// DeleteReservationHandler maneja la eliminación de una reserva
func DeleteReservationHandler(reservationService service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, "Invalid reservation ID", err.Error())
		}

		if err := reservationService.DeleteReservation(uint(id)); err != nil {
			log.Debug("Error deleting reservation:", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "Could not delete reservation", err.Error())
		}

		return utils.HandleResponse(c, fiber.StatusOK, "Reservation deleted successfully", nil)
	}
}
