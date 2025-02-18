package service

import (
	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
)

type PaymentService interface {
	CreatePayment(payment *models.Payment) error
	GetPaymentByID(id uint, params map[string]interface{}) (*models.Payment, error)
	GetAllPayments(page, pageSize int, id int, preload bool) ([]models.Payment, int64, error)
	UpdatePayment(payment *models.Payment) error
	DeletePayment(id uint) error
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentService{paymentRepo}
}

func (s *paymentService) CreatePayment(payment *models.Payment) error {

	return s.paymentRepo.Create(payment)
}

func (s *paymentService) GetPaymentByID(id uint, params map[string]interface{}) (*models.Payment, error) {

	return s.paymentRepo.FindByID(id)
}

func (s *paymentService) GetAllPayments(page, pageSize int, id int, preload bool) ([]models.Payment, int64, error) {
	//var apartments []models.Apartment
	//var err error
	//var totalRecords int64
	/*
	   buildingId := params["id"].(int)

	   if buildingId != 0 {

	   		apartments, totalRecords, err = s.paymentRepo.FindAllByBuildingID(page, pageSize, "Residents", uint(buildingId))
	   		if err != nil {
	   			return nil, 0, err
	   		}
	   		return apartments, totalRecords, nil
	   	}
	*/return nil, 0, nil
}

func (s *paymentService) UpdatePayment(payment *models.Payment) error {
	payment, err := s.paymentRepo.FindByID(payment.ID)
	if err != nil {
		return err
	}

	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}
	return nil
}

func (s *paymentService) DeletePayment(id uint) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.paymentRepo.Delete(payment)
}
