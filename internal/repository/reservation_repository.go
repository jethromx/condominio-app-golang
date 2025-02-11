package repository

import (
	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(reservation *models.Reservation) error
	FindByID(id uint) (*models.Reservation, error)
	FindAll() ([]models.Reservation, error)
	Update(reservation *models.Reservation) error
	Delete(reservation *models.Reservation) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db}
}

func (r *reservationRepository) Create(reservation *models.Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.First(&reservation, id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) Update(reservation *models.Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *reservationRepository) Delete(reservation *models.Reservation) error {
	return r.db.Delete(reservation).Error
}
