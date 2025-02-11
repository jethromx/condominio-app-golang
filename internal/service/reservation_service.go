package service

import (
	"time"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
)

type ReservationService interface {
	CreateReservation(startDate, endDate time.Time, residentID uint) (*models.Reservation, error)
	GetReservationByID(id uint) (*models.Reservation, error)
	GetAllReservations() ([]models.Reservation, error)
	UpdateReservation(id uint, startDate, endDate time.Time, residentID uint) (*models.Reservation, error)
	DeleteReservation(id uint) error
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
}

func NewReservationService(reservationRepo repository.ReservationRepository) ReservationService {
	return &reservationService{reservationRepo}
}

func (s *reservationService) CreateReservation(startDate, endDate time.Time, residentID uint) (*models.Reservation, error) {
	reservation := &models.Reservation{
		// StartDate:  startDate,
		// EndDate:    endDate,
		ResidentID: residentID,
	}
	if err := s.reservationRepo.Create(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func (s *reservationService) GetReservationByID(id uint) (*models.Reservation, error) {
	return s.reservationRepo.FindByID(id)
}

func (s *reservationService) GetAllReservations() ([]models.Reservation, error) {
	return s.reservationRepo.FindAll()
}

func (s *reservationService) UpdateReservation(id uint, startDate, endDate time.Time, residentID uint) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	//reservation.StartDate = startDate
	// reservation.EndDate = endDate
	reservation.ResidentID = residentID
	if err := s.reservationRepo.Update(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func (s *reservationService) DeleteReservation(id uint) error {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.reservationRepo.Delete(reservation)
}
