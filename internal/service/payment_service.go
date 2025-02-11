package service

import (
	"time"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
)

type PaymentService interface {
	CreatePayment(amount float64, date time.Time, residentID uint) (*models.Payment, error)
	GetPaymentByID(id uint) (*models.Payment, error)
	GetAllPayments() ([]models.Payment, error)
	UpdatePayment(id uint, amount float64, date time.Time, residentID uint) (*models.Payment, error)
	DeletePayment(id uint) error
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentService{paymentRepo}
}

func (s *paymentService) CreatePayment(amount float64, date time.Time, residentID uint) (*models.Payment, error) {
	payment := &models.Payment{
		Amount: amount,
		//Date:       date,
		ResidentID: residentID,
	}
	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *paymentService) GetPaymentByID(id uint) (*models.Payment, error) {
	return s.paymentRepo.FindByID(id)
}

func (s *paymentService) GetAllPayments() ([]models.Payment, error) {
	return s.paymentRepo.FindAll()
}

func (s *paymentService) UpdatePayment(id uint, amount float64, date time.Time, residentID uint) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	payment.Amount = amount
	//payment.Date = date
	payment.ResidentID = residentID
	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *paymentService) DeletePayment(id uint) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.paymentRepo.Delete(payment)
}
